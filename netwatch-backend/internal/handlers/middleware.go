package handlers

import (
	"strings"

	"github.com/ertel/netwatch-backend/internal/models"
	"github.com/ertel/netwatch-backend/internal/services"
	"github.com/ertel/netwatch-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware valida o JWT e injeta os claims no contexto
func AuthMiddleware(authSvc services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(c, "token não fornecido")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.Unauthorized(c, "formato de token inválido")
		}

		claims, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			return utils.Unauthorized(c, "token inválido ou expirado")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// RequireRole garante que o usuário tem a role mínima necessária
func RequireRole(minRole models.Role) fiber.Handler {
	roleLevel := map[models.Role]int{
		models.RoleViewer:   1,
		models.RoleOperator: 2,
		models.RoleAdmin:    3,
	}

	return func(c *fiber.Ctx) error {
		roleStr, ok := c.Locals("userRole").(string)
		if !ok {
			return utils.Unauthorized(c, "")
		}

		userLevel := roleLevel[models.Role(roleStr)]
		requiredLevel := roleLevel[minRole]

		if userLevel < requiredLevel {
			return utils.Forbidden(c, "permissão insuficiente")
		}

		return c.Next()
	}
}

// ParsePage extrai parâmetros de paginação da query string
func ParsePage(c *fiber.Ctx) (page, limit int) {
	page = c.QueryInt("page", 1)
	limit = c.QueryInt("limit", 50)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	return page, limit
}
