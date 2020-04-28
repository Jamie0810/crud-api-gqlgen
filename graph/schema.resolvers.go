package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/jamie/gqlgen-crud/graph/generated"
	"github.com/jamie/gqlgen-crud/graph/model"
	"github.com/jamie/gqlgen-crud/models"
	mysql "github.com/jamie/gqlgen-crud/mysql"
)

var db = mysql.Connection()

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*models.Todo, error) {
	todo := models.Todo{
		Text:   input.Text,
		UserID: input.UserID,
		Done:   false,
	}
	db.Create(&todo)
	return &todo, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input model.EditTodo) (*models.Todo, error) {
	todo := models.Todo{ID: input.ID}
	db.First(&todo)
	todo.Text = input.Text
	db.Model(&models.Todo{}).Update(&todo)
	return &todo, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, input int) (*models.Todo, error) {
	todo := models.Todo{
		ID: input,
	}
	db.First(&todo)
	db.Delete(&todo)
	return &todo, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*models.User, error) {
	user := models.User{
		Name: input.Name,
	}
	db.Create(&user)
	return &user, nil
}

func (r *queryResolver) Todos(ctx context.Context) ([]*models.Todo, error) {
	var todos []*models.Todo
	db.Preload("User").Find(&todos)
	return todos, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	db.Find(&users)
	return users, nil
}

func (r *queryResolver) Todo(ctx context.Context, input *model.FetchTodo) (*models.Todo, error) {
	var todo models.Todo
	db.Preload("User").First(&todo, input.ID)
	return &todo, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
