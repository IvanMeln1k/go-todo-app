package handler

import "github.com/labstack/echo/v4"

func (h *Handler) signUp(c echo.Context) error {
	return c.String(200, c.Path());
}

func (h *Handler) signIn(c echo.Context) error {
	return c.String(200, c.Path())
}