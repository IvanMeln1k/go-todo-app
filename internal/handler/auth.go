package handler

import (
	"net/http"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (h *Handler) signUp(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(user); err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}
	id, err := h.services.CreateUser(*user);
	if err != nil {
		logrus.Errorf("%s", err);
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(200, map[string]interface{}{
		"id": id,
	});
}

func (h *Handler) signIn(c echo.Context) error {
	return c.String(200, c.Path())
}