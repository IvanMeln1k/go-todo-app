package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (h *Handler) createItem(c echo.Context) error {
	return c.String(200, c.Path());
}

func (h *Handler) getAllItems(c echo.Context) error {
	fmt.Println("hihihi")
	return c.String(200, c.Path())
}

func (h *Handler) getItemById(c echo.Context) error {
	return c.String(200, c.Path())
}

func (h *Handler) updateItem(c echo.Context) error {
	return c.String(200, c.Path())
}

func (h *Handler) deleteItem(c echo.Context) error {
	return c.String(200, c.Path())
}