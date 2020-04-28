# Create a CRUD API with gqlgen
In this walkthrough, we will be creating a GraphQL API in Golang with MySQL and GORM using gqlgen.

## What is [gqlgen](https://github.com/99designs/gqlgen) ?

It’s a Golang library which allows us to auto generate Golang code stub from [GraphQL](https://graphql.org/) schema automatically.

## Setup Project
```
mkdir gqlgen-crud
cd gqlgen-crud
go mod init github.com/jamie/gqlgen-crud
go get github.com/99designs/gqlgen
gqlgen init
```
`gqlgen init` will initialize the **gqlgen server**. Now once it's done, your project structure should look like this:
```
.
├── go.mod
├── go.sum
├── gqlgen.yml
├── graph
│   ├── generated
│   │   └── generated.go
│   ├── model
│   │   └── models_gen.go
│   ├── resolver.go  
│   ├── schema.graphqls
│   ├── schema.resolvers.go
└── server.go
```

[Getting Started tutorial](https://gqlgen.com/getting-started/) demonstrates how to build an API service with [GraphQL](https://graphql.org/learn/schema/) by using gqlgen. The model in this tutorial is able to "**Create**" and "**Query**" a to-do list. Now we are going to add "**Update**" and "**Delete**".

***Note**
If you would like to make modifications to fit the project layout, you will only have to edit two files **before regenerating the project** : 
1. **`gqlgen.yml`** - To direct gqlgen about where to auto generate files 
1. **`schema.graphql`** - To describe our API using the schema definition library.

---
## Modify GraphQL schema & yml 
First, we create a new GraphQL schema file `schema.graphql` at root directory:
```schema
type Todo {
  id:   Int!
  text: String!
  done: Boolean!
  userID: Int!
  user: User!
}

type User {
  id: Int!
  name: String!
}

type Query {
  todos: [Todo!]!
  users: [User!]!
  todo(input: FetchTodo): Todo!
}

input NewTodo {
  text: String!
  userId: Int!
}

input EditTodo {
  id: Int!
  text: String!
}

input NewUser {
  name: String!
}

input FetchTodo {
  id: Int!
}

type Mutation {
  createTodo(input: NewTodo!): Todo!
  updateTodo(input: EditTodo!): Todo!
  deleteTodo(input: Int!): Todo!
  createUser(input: NewUser!): User!
}
```
Here we have defined our basic models and some mutations to create, update and delete the to-do list, and queries to get all users and todos.

Then modify the `gqlgen.yml`:

```yaml
schema:
    - schema.graphql

exec:
    filename: graph/generated/generated.go
    package: generated

model:
    filename: graph/model/models_gen.go

resolver:
    layout: follow-schema
    dir: graph
    package: graph
```

## Prepare Gorm Model
Now we are going to define the model for database schema. Since we are using [GORM](https://gorm.io/docs/) as our ORM, we would need to add GORM-specific tags to our **Todo** and **User** struct. Hence it makes sense to create our own struct, and instruct gqlgen to use them instead of the auto-generated ones.


Let’s create `models/models.go` and place our model structs for Todo and User:

```go
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
```

Next, we need to instruct gqlgen to use this by adding an entry to the `gqlgen.yml` file to refer to the new models we created:
```yaml
models:
    Todo:    
        model: github.com/jamie/gqlgen-crud/models.Todo
    User:
        model: github.com/jamie/gqlgen-crud/models.User
```

Then **regenerate the project**, run:
```
gqlgen
```

Now the `models_gen.go` should look like this:
```go
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
```

## Initial DB setup and configuration
Here we create the database and tables in `mysql/connection.go` with GORM:

```go
package database

import (
    "fmt"

    _ "github.com/go-sql-driver/mysql"
    "github.com/jamie/gqlgen-crud/models"
    "github.com/jinzhu/gorm"
)

func Connection() *gorm.DB {
    db, err := gorm.Open("mysql", "USER:PASSWORD@(localhost)/")

    if err != nil {
        fmt.Println(err)
        panic("Failed to connect to database")
    }
    
    
    // Create the database
    db.Exec("CREATE DATABASE IF NOT EXISTS db_name")
    db.Exec("USE db_name")
    
    //Migration to create tables for Todo and User schema
    db.AutoMigrate(&models.Todo{}, &models.User{}) 


    return db
}
```

Remember to import the set of package in `schema.resolvers.go`:
```go
import (
    mysql "github.com/jamie/gqlgen-crud/mysql"
)

var db = mysql.Connection()
```

## Implemntaion Resolver
Last thing!! We need to implement the real logic of the interface method. 

Open `schema.resolvers.go` and provide the definition to mutation and queries. The stubs are already auto-generated **with a not implemented panic statement** , so let’s override that.
```go
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
```
---
## Running & Testing the App

Run the `server.go` file. The server will start at [localhost:8080](localhost:8080)
```
go run server.go
```
Now you can try the mutations or queries below:
```graphql
mutation createTodo {
  createTodo(input:{text: "todo3", userId: 2}) {
    id
    text
    done
    userID
    user{
      id
      name
    }
  }
}

mutation createUser {
  createUser(input:{name: "Jamie"}) {
    id
    name
  }
}

mutation updateTodo {
  updateTodo(input:{id: 1, text: "update"}) {
    id
    text
    done
    userID
    user{
      id
      name
    }
  }
}

mutation deleteTodo {
  deleteTodo(input: 1) {
    text
    done
  }
}


query findTodos {
  	todos {
      id
      text
      done
      userID
      user {
        id
        name
      }
    }
}

query findTodo {
  	todo(input: { id: 2}) {
    id
      text
      done
      userID
      user {
        id
        name
      }
    }
}

query findUsers {
  	users {
      id
      name
    }
}
```
