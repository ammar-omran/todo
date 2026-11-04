package app

import (
	"html/template"
	"net/http"
	"to-do/api"
	"to-do/handler"
	"to-do/views"
)

func (a *App) loadRoutes(tmpl *template.Template) {
	todoHandler := handler.New(a.logger, a.dbConn)
	todoApi := &api.Api{Th: todoHandler}
	homeView := &views.Home{Th: todoHandler, Tmpl: tmpl}

	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("GET /", todoApi.Get)
	apiRouter.HandleFunc("POST /", todoApi.Create)
	apiRouter.HandleFunc("PUT /{id}", todoApi.Update)
	apiRouter.HandleFunc("DELETE /{id}", todoApi.Delete)

	webRouter := http.NewServeMux()
	files := http.FileServer(http.Dir("./static"))

	webRouter.Handle(
		"GET /static/",
		http.StripPrefix("/static/", files),
	)
	webRouter.HandleFunc("GET /", homeView.Home)
	webRouter.HandleFunc("GET /todos", homeView.List)
	webRouter.HandleFunc("POST /todos", homeView.Create)
	webRouter.HandleFunc("PUT /todos/{id}", homeView.Update)
	webRouter.HandleFunc("DELETE /todos/{id}", homeView.Delete)

	a.router.Handle("/api/todos/", http.StripPrefix("/api/todos", apiRouter))
	a.router.Handle("/", webRouter)
}
