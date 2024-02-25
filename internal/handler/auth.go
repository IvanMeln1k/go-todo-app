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
	id, err := h.services.CreateUser(*user)
	if err != nil {
		logrus.Errorf("%s", err)
		if err.Error() == "username already in use" {
			return newErrorResponse(409, "Username already in use")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(200, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *Handler) signIn(c echo.Context) error {
	user := new(signInInput)
	if err := c.Bind(user); err != nil {
		return newErrorResponse(400, err.Error())
	}
	if err := c.Validate(user); err != nil {
		return newErrorResponse(400, err.Error())
	}
	token, err := h.services.Authorization.GenerateToken(user.Username, user.Password)
	if err != nil {
		if err.Error() == "user not found" {
			return newErrorResponse(401, "Invalid username or password")
		}
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, map[string]interface{}{
		"token": token,
	})
}
