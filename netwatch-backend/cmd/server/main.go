package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ertel/netwatch-backend/internal/config"
	"github.com/ertel/netwatch-backend/internal/database"
	"github.com/ertel/netwatch-backend/internal/handlers"
	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/repository"
	"github.com/ertel/netwatch-backend/internal/services"
	"github.com/ertel/netwatch-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

func main() {
	// ----------------------------------------------------------------
	// 1. Configuração
	// ----------------------------------------------------------------
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: falha ao carregar configuração: %v\n", err)
		os.Exit(1)
	}

	// ----------------------------------------------------------------
	// 2. Logger (Zap)
	// ----------------------------------------------------------------
	log := buildLogger(cfg)
	defer log.Sync() //nolint:errcheck

	log.Info("Iniciando NetWatch Backend",
		zap.String("env", cfg.Server.Env),
		zap.Int("port", cfg.Server.Port),
	)

	// ----------------------------------------------------------------
	// 3. Banco de Dados
	// ----------------------------------------------------------------
	db, err := database.Connect(cfg, log)
	if err != nil {
		log.Fatal("Falha ao conectar no banco de dados", zap.Error(err))
	}
	defer database.Close(db) //nolint:errcheck

	// Executar migrations
	if err := database.RunMigrations(db, log); err != nil {
		log.Fatal("Falha ao executar migrations", zap.Error(err))
	}

	// ----------------------------------------------------------------
	// 4. Dependências (Repositories, Services, Handlers)
	// ----------------------------------------------------------------
	userRepo   := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	metricRepo := repository.NewMetricRepository(db)

	authSvc   := services.NewAuthService(userRepo, cfg)
	deviceSvc := services.NewDeviceService(deviceRepo, log)

	authHandler   := handlers.NewAuthHandler(authSvc, log)
	deviceHandler := handlers.NewDeviceHandler(deviceSvc, log)
	metricHandler := handlers.NewMetricHandler(metricRepo, log)

	// ----------------------------------------------------------------
	// 5. Fiber App
	// ----------------------------------------------------------------
	app := fiber.New(fiber.Config{
		AppName:      "NetWatch API v0.1.0",
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Error("erro não tratado", zap.Error(err), zap.String("path", c.Path()))
			return utils.InternalError(c, "erro interno do servidor")
		},
	})

	// Middlewares globais
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Join(cfg.CORS.AllowedOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimit.Requests,
		Expiration: time.Duration(cfg.RateLimit.WindowSeconds) * time.Second,
	}))

	// ----------------------------------------------------------------
	// 6. Rotas
	// ----------------------------------------------------------------
	api := app.Group("/api")

	// Health check (público)
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := database.Ping(db); err != nil {
			return utils.InternalError(c, "banco de dados inacessível")
		}
		return utils.OK(c, fiber.Map{
			"status":  "ok",
			"version": "0.1.0",
			"time":    time.Now().UTC(),
		})
	})

	// Auth (público)
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	// Rotas autenticadas
	protected := api.Group("/", handlers.AuthMiddleware(authSvc))

	protected.Post("/auth/logout", authHandler.Logout)
	protected.Get("/auth/me", authHandler.Me)

	// Admin: gestão de usuários
	adminAuth := protected.Group("/auth", handlers.RequireRole(models.RoleAdmin))
	adminAuth.Post("/register", authHandler.Register)

	// Devices (operator+)
	devices := protected.Group("/devices")
	devices.Get("/", deviceHandler.List)
	devices.Get("/:id", deviceHandler.Get)
	devices.Post("/", handlers.RequireRole(models.RoleOperator), deviceHandler.Create)
	devices.Put("/:id", handlers.RequireRole(models.RoleOperator), deviceHandler.Update)
	devices.Delete("/:id", handlers.RequireRole(models.RoleAdmin), deviceHandler.Delete)

	// Metrics (viewer+)
	metrics := protected.Group("/metrics")
	metrics.Get("/:deviceId/history", metricHandler.GetHistory)
	metrics.Get("/:deviceId/summary", metricHandler.GetSummary)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return utils.NotFound(c, "rota não encontrada")
	})

	// ----------------------------------------------------------------
	// 7. Iniciar servidor com graceful shutdown
	// ----------------------------------------------------------------
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Servidor HTTP iniciado", zap.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			log.Fatal("Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	<-quit
	log.Info("Encerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("Erro no graceful shutdown", zap.Error(err))
	}

	log.Info("Servidor encerrado com sucesso")
}

func buildLogger(cfg *config.Config) *zap.Logger {
	level := zap.InfoLevel
	switch cfg.Logging.Level {
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	zapCfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      !cfg.IsProduction(),
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if !cfg.IsProduction() {
		zapCfg.Encoding = "console"
	}

	log, err := zapCfg.Build()
	if err != nil {
		panic(fmt.Sprintf("falha ao construir logger: %v", err))
	}

	return log
}
