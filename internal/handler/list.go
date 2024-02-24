package handler

import (
	"strconv"

	"github.com/IvanMeln1k/go-todo-app/internal/domain"
	"github.com/labstack/echo/v4"
)

func (h *Handler) createList(c echo.Context) error {
	userId, err := getUserId(c)

	if err != nil {
		return err
	}

	todoList := new(domain.TodoList)
	if err = c.Bind(todoList); err != nil {
		return newErrorResponse(400, err.Error())
	}
	if err = c.Validate(todoList); err != nil {
		return newErrorResponse(400, err.Error())
	}

	todoListId, err := h.services.TodoList.Create(userId, *todoList)
	if err != nil {
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(201, todoListId)
}

func (h *Handler) getAllLists(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	todoLists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		return newErrorResponse(404, err.Error())
	}

	return c.JSON(200, map[string]interface{}{
		"todoLists": todoLists,
	})
}

func (h *Handler) getListById(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	todoListId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(400, "Bad request")
	}

	todoList, err := h.services.TodoList.GetById(userId, todoListId)
	if err != nil {
		if err.Error() == "not found" {
			return newErrorResponse(404, "Not found")
		}
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(200, map[string]interface{}{
		"todoList": todoList,
	})
}

func (h *Handler) updateList(c echo.Context) error {
	return c.String(200, c.Path())
}

func (h *Handler) deleteList(c echo.Context) error {
	return c.String(200, c.Path())
}