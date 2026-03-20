package utils

import "github.com/gofiber/fiber/v2"

// APIResponse é a estrutura padrão de resposta da API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
}

// PaginationMeta contém informações de paginação
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// OK envia resposta 200 com dados
func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Data:    data,
	})
}

// OKPaginated envia resposta 200 com dados paginados
func OKPaginated(c *fiber.Ctx, data interface{}, page, limit int, total int64) error {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Success: true,
		Data:    data,
		Meta: &PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Created envia resposta 201
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent envia resposta 204
func NoContent(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNoContent).Send(nil)
}

// BadRequest envia resposta 400
func BadRequest(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(APIResponse{
		Success: false,
		Error:   msg,
	})
}

// Unauthorized envia resposta 401
func Unauthorized(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "não autorizado"
	}
	return c.Status(fiber.StatusUnauthorized).JSON(APIResponse{
		Success: false,
		Error:   msg,
	})
}

// Forbidden envia resposta 403
func Forbidden(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "acesso negado"
	}
	return c.Status(fiber.StatusForbidden).JSON(APIResponse{
		Success: false,
		Error:   msg,
	})
}

// NotFound envia resposta 404
func NotFound(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "recurso não encontrado"
	}
	return c.Status(fiber.StatusNotFound).JSON(APIResponse{
		Success: false,
		Error:   msg,
	})
}

// Conflict envia resposta 409
func Conflict(c *fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusConflict).JSON(APIResponse{
		Success: false,
		Error:   msg,
	})
}

// InternalError envia resposta 500
func InternalError(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "erro interno do servidor"
	}
	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{
		Success: false,
		Error:   msg,
	})
}
