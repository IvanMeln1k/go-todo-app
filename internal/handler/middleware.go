package handler

import (
	"errors"
	"strings"

	"github.com/IvanMeln1k/go-todo-app/internal/service"
	"github.com/labstack/echo/v4"
)

func (h *Handler) userIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return newErrorResponse(401, "Unauthorized")
		}
		params := strings.Split(authHeader, " ")
		if len(params) < 2 {
			return newErrorResponse(401, "Unauthorized")
		}
		if params[0] != "Bearer" {
			return newErrorResponse(401, "Unauthorized")
		}

		userId, err := h.services.Authorization.ParseToken(params[1])

		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				return newErrorResponse(401, "Token is expired")
			} else if errors.Is(err, service.ErrInvalidTokenSignature) {
				return newErrorResponse(401, "Invalid token signature")
			}
			return newErrorResponse(500, "Internal server error")
		}

		c.Set("userId", userId)

		return next(c)
	}
}

func getUserId(c echo.Context) (int, error) {
	id := c.Get("userId")

	idInt, ok := id.(int)
	if !ok {
		return 0, newErrorResponse(401, "Unautharized")
	}

	return idInt, nil
}
