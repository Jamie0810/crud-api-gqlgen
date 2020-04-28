package models

type Todo struct {
	ID     int
	Text   string
	Done   bool
	UserID int
	User   User
}

type User struct {
	ID    int
	Name  string
	Todos []Todo
}
