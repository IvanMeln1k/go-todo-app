package domain

type TodoList struct {
	Id int `json:"id" db:"id"`
	Title string `json:"title" validate:"required" db:"title"`
	Description string `json:"description" title:"description"`
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