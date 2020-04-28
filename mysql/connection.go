package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jamie/gqlgen-crud/models"
	"github.com/jinzhu/gorm"
)

func Connection() *gorm.DB {
	db, err := gorm.Open("mysql", "root:Foxconn123@(localhost)/")

	if err != nil {
		fmt.Println(err)
		panic("Failed to connect to database")
	}
	db.Exec("CREATE DATABASE IF NOT EXISTS gqlgen")
	db.Exec("USE gqlgen")
	db.AutoMigrate(&models.Todo{}, &models.User{})

	return db
}
