package handler

import (
	"errors"
	"net/http"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/IvanMeln1k/go-todo-app/internal/service"
	"github.com/labstack/echo/v4"
)

func (h *Handler) signUp(c echo.Context) error {
	user := new(domain.User)
	if err := c.Bind(user); err != nil {
		return newErrorResponse(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(user); err != nil {
		return newErrorResponse(http.StatusBadRequest, "invalid body")
	}
	id, err := h.services.CreateUser(*user)
	if err != nil {
		if errors.Is(err, service.ErrUsernameAlreadyInUse) {
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

	tokens, err := h.services.Authorization.SignIn(c.Request().Context(), user.Username, user.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return newErrorResponse(401, "Invalid username or password")
		} else if errors.Is(err, service.ErrInternal) {
			return newErrorResponse(500, "Internal server error")
		}
		return newErrorResponse(401, "Anauthorized")
	}

	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})
	return c.JSON(200, map[string]interface{}{
		"tokens": tokens,
	})
}

func (h *Handler) refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return newErrorResponse(401, "Unauthorized")
	}
	tokens, err := h.services.Authorization.Refresh(c.Request().Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrSessionExpiredOrInvalid) {
			return newErrorResponse(401, "Unauthorized")
		} else if errors.Is(err, service.ErrInvalidSession) {
			return newErrorResponse(401, "Invalid session")
		} else if errors.Is(err, service.ErrInternal) {
			return newErrorResponse(500, "Internal server error")
		}
		return newErrorResponse(401, "Unauthorized")
	}
	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
	})
	return c.JSON(200, map[string]interface{}{
		"tokens": tokens,
	})
}

func (h *Handler) logout(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return newErrorResponse(401, "Unauthorized")
	}
	err = h.services.Authorization.Logout(c.Request().Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrSessionExpiredOrInvalid) {
			return newErrorResponse(401, "Unauthorized")
		}
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) logoutAll(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return newErrorResponse(401, "Unauthorized")
	}
	err = h.services.LogoutAll(c.Request().Context(), refreshToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrSessionExpiredOrInvalid) {
			return newErrorResponse(401, "Unauthorized")
		}
		return newErrorResponse(500, "Internal server error")
	}
	return c.JSON(200, map[string]interface{}{
		"status": "ok",
	})
}
