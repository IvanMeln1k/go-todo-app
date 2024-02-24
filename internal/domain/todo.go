package domain

import "errors"

type TodoList struct {
	Id int `json:"id" db:"id"`
	Title string `json:"title" validate:"required" db:"title"`
	Description string `json:"description" title:"description"`
}

type UpdateTodoList struct {
	Title *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateTodoList) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update struct has no values")
	}
	return nil
}

type UsersList struct {
	Id int
	UserId int
	ListId int
}

type TodoItem struct {
	Id int `json:"id"`
	Title string `json:"title"`
 	Description string `json:"description"`
	Done bool `done:"done"`
}

type ListsItem struct {
	Id int
	ListId int
	ItemId int
}