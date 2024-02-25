package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func newErrorResponse(statusCode int, message string) error {
	logrus.Error(message)
	return echo.NewHTTPError(statusCode, message)
}
