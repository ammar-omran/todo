package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"to-do/data"

	_ "github.com/mattn/go-sqlite3"
)

type Repository[T struct{}] interface {
	Add() error
	Delete() error
	Complete() error
	List(bool) []T
}

func main() {
	conn, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	db := data.New(conn)

	if _, err := db.CreateToDo(
		context.Background(),
		func(d data.ToDo) data.CreateToDoParams {
			return data.CreateToDoParams{
				Description: d.Description,
				Created:     d.Created,
				Done:        d.Done,
			}
		}(data.CreateToDo("Hello, World")),
	); err != nil {
		fmt.Println(err)
	}

	if res, err := db.GetToDosByState(context.Background(), false); err == nil {
		fmt.Println(res)
	} else {
		fmt.Println(err)
	}

	if err := db.ChangeToDoStatus(context.Background(), 50, false); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Done Updating")
	}
}
