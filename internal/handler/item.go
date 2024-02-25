package handler

import (
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

	return c.JSON(201, todoItemId)
}

func (h *Handler) getAllItems(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	todoListId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(400, "TodoListId is no integer value")
	}

	todoItems, err := h.services.TodoItem.GetAll(userId, todoListId)
	if err != nil {
		if err.Error() == "not found" {
			return newErrorResponse(404, "Not found")
		}
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(200, map[string]interface{}{
		"todoItems": todoItems,
	})
}

func (h *Handler) getItemById(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	todoItemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(400, "TodoItemId is no integer value")
	}

	todoItem, err := h.services.TodoItem.GetById(userId, todoItemId)
	if err != nil {
		if err.Error() == "not found" {
			return newErrorResponse(404, "Item not found")
		}
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(200, map[string]interface{}{
		"todoItem": todoItem,
	})
}

func (h *Handler) deleteItem(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	todoItemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(400, "TodoItemId is no integer value")
	}

	err = h.services.TodoItem.Delete(userId, todoItemId)
	if err != nil {
		if err.Error() == "not found" {
			return newErrorResponse(404, "Item not found")
		}
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(200, map[string]interface{}{
		"status": "ok",
	})
}

func (h *Handler) updateItem(c echo.Context) error {
	userId, err := getUserId(c)
	if err != nil {
		return err
	}

	todoItemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return newErrorResponse(400, "TodoItemId is no integer value")
	}

	var updateTodoItem domain.UpdateTodoItem
	if err = c.Bind(&updateTodoItem); err != nil {
		return newErrorResponse(400, err.Error())
	}
	if err = updateTodoItem.Validate(); err != nil {
		return newErrorResponse(400, err.Error())
	}

	todoItem, err := h.services.TodoItem.Update(userId, todoItemId, updateTodoItem)
	if err != nil {
		if err.Error() == "not found" {
			return newErrorResponse(404, "TodoItem not found")
		}
		return newErrorResponse(500, "Internal server error")
	}

	return c.JSON(201, map[string]interface{}{
		"todoItem": todoItem,
	})
}
