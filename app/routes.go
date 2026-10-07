package app

import (
	"html/template"
	"net/http"
	"to-do/api"
	"to-do/handler"
)

func (a *App) loadRoutes(tmpl *template.Template) {
	todoHandler := handler.New(a.logger, a.dbConn)
	todoApi := &api.Api{Th: todoHandler}

	apiRouter := http.NewServeMux()

	apiRouter.HandleFunc("GET /", todoApi.Get)
	apiRouter.HandleFunc("POST /", todoApi.Create)
	apiRouter.HandleFunc("PUT /{id}", todoApi.Update)
	apiRouter.HandleFunc("DELETE /{id}", todoApi.Delete)

	a.router.Handle("/api/todos/", http.StripPrefix("/api/todos", apiRouter))
}
