package handler

import "github.com/labstack/echo/v4"

func (h *Handler) createList(c echo.Context) error {
	return c.String(200, c.Path());
}

func (h *Handler) getAllLists(c echo.Context) error {
	return c.String(200, c.Path())
}

func (h *Handler) getListById(c echo.Context) error {
	return c.String(200, c.Path())
}

func (h *Handler) updateList(c echo.Context) error {
	return c.String(200, c.Path())
}

func (h *Handler) deleteList(c echo.Context) error {
	return c.String(200, c.Path())
}