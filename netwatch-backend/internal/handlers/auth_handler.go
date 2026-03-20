package handlers

import (
	"errors"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/services"
	"github.com/ertel/netwatch-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authSvc services.AuthService
	log     *zap.Logger
}

func NewAuthHandler(authSvc services.AuthService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, log: log}
}

// Login godoc
// POST /api/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "corpo da requisição inválido")
	}

	tokens, user, err := h.authSvc.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			return utils.Unauthorized(c, "email ou senha incorretos")
		case errors.Is(err, services.ErrUserInactive):
			return utils.Unauthorized(c, "usuário desativado")
		default:
			h.log.Error("erro no login", zap.Error(err))
			return utils.InternalError(c, "")
		}
	}

	return utils.OK(c, fiber.Map{
		"tokens": tokens,
		"user":   user.ToResponse(),
	})
}

// Logout godoc
// POST /api/auth/logout
// (stateless JWT - client descarta o token)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return utils.OK(c, fiber.Map{"message": "logout realizado com sucesso"})
}

// RefreshToken godoc
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
		return utils.BadRequest(c, "refresh_token é obrigatório")
	}

	tokens, err := h.authSvc.RefreshToken(body.RefreshToken)
	if err != nil {
		return utils.Unauthorized(c, "refresh token inválido ou expirado")
	}

	return utils.OK(c, tokens)
}

// Me godoc
// GET /api/auth/me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	return utils.OK(c, fiber.Map{
		"id":    c.Locals("userID"),
		"email": c.Locals("userEmail"),
		"role":  c.Locals("userRole"),
	})
}

// Register godoc
// POST /api/auth/register (admin only)
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "corpo da requisição inválido")
	}

	user, err := h.authSvc.Register(req)
	if err != nil {
		h.log.Error("erro ao registrar usuário", zap.Error(err))
		return utils.InternalError(c, "erro ao criar usuário")
	}

	return utils.Created(c, user.ToResponse())
}
