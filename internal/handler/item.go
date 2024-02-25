package handler

import (
	"fmt"
	"strconv"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/labstack/echo/v4"
)

func (h *Handler) createItem(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	todoListId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(400, "TodoList is no integer value")
	}

	var todoItem domain.TodoItem
	if err = c.Bind(&todoItem); err != nil {
		return newErrorResponse(400, err.Error())
	}
	if err = c.Validate(&todoItem); err != nil {
		return newErrorResponse(400, err.Error())
	}

	todoItemId, err := h.services.TodoItem.Create(userId, todoListId, todoItem)
	if err != nil {
		if err.Error() == "not found" {
			return newErrorResponse(404, "Not found")
		}
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(201, todoItemId);
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