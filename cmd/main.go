package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"to-do/api"
	"to-do/data"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

type Repository[T struct{}] interface {
	Add() error
	Delete() error
	Complete() error
	List(bool) []T
}

var db *data.Queries

func create(w http.ResponseWriter, req *http.Request) {
	var desc = req.FormValue("description")
	if desc == "" {
		api.RespondWithError(w, 400, "Description can't be empty!")
	}
	if _, err := db.CreateToDo(
		context.Background(),
		func(d data.ToDo) data.CreateToDoParams {
			return data.CreateToDoParams{
				Description: d.Description,
				Created:     d.Created,
				Done:        d.Done,
			}
		}(data.CreateToDo(desc)),
	); err != nil {
		api.RespondWithError(w, 422, err.Error())
	}

	api.RespondWithJSON(w, 201, struct {
		Status  bool
		Message string
	}{true, "done"})

}

func update(w http.ResponseWriter, req *http.Request) {
	var status = req.URL.Query().Get("status")
	completed := status != "reopen"
	if err := db.ChangeToDoStatus(context.Background(), 50, completed); err != nil {
		api.RespondWithError(w, 500, err.Error())
	} else {
		api.RespondWithJSON(w, 201, struct {
			Status  bool
			Message string
		}{true, "Done"})
	}
}

func get(w http.ResponseWriter, req *http.Request) {
	var status = req.URL.Query().Get("status")
	completed := status == "completed"
	if res, err := db.GetToDosByState(context.Background(), completed); err == nil {
		api.RespondWithJSON(w, 200, struct {
			Status  bool
			Message string
			Data    any
		}{true, "fetched", res})
	} else {
		api.RespondWithError(w, 500, err.Error())
	}
}

func delete(w http.ResponseWriter, req *http.Request) {
	var desc = req.FormValue("description")
	if desc == "" {
		api.RespondWithError(w, 400, "Description can't be empty!")
	}
	if res, err := db.GetToDosByState(context.Background(), false); err == nil {
		fmt.Println(res)
	} else {
		fmt.Println(err)
	}
}

func main() {
	conn, err := sql.Open("sqlite", "todo.db")
	if err != nil {
		log.Fatal(err)
	}
	db = data.New(conn)

	http.HandleFunc("POST /create", create)
	http.HandleFunc("GET /get", get)
	http.HandleFunc("PATCH /set-status", update)
	http.HandleFunc("DELETE /delete", delete)
	if err = http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}

}

func NewServeCommand(showStartBanner bool) *cobra.Command {

	command := &cobra.Command{
		Use:          "serve",
		Args:         cobra.ArbitraryArgs,
		Short:        "Starts the web server (default to 127.0.0.1:8090)",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			err := http.ListenAndServe(":8090", nil)
			return err
		},
	}

	return command
}
