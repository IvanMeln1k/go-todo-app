package handler

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) userIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func (c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		
		fmt.Println(authHeader)
		
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
			if err.Error() == "Token is expired" {
				return newErrorResponse(401, "Token is expired")
			} else if err.Error() == "invalid token signature" {
				return newErrorResponse(401, "Ivanlid token signature")
			}
			return newErrorResponse(500, "Internal server error")
		}

		c.Set("userId", userId)

		return next(c)
	}
}