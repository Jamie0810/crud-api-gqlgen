// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type EditTodo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type FetchTodo struct {
	ID int `json:"id"`
}

type NewTodo struct {
	Text   string `json:"text"`
	UserID int    `json:"userId"`
}

type NewUser struct {
	Name string `json:"name"`
}
