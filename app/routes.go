package app

import (
	"html/template"
	"net/http"
	"os"
	"to-do/api"
	"to-do/handler"
	"to-do/views"

	"ganigeorgiev.com/pocketbase/pocketbase/apis"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

func (a *App) loadRoutes(tmpl *template.Template) {
	todoHandler := handler.New(a.logger, a.pb)
	todoApi := &api.Api{Th: todoHandler}
	homeView := &views.Home{Th: todoHandler, Tmpl: tmpl}

	a.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./static"), false))

		se.Router.GET("/", apis.WrapStdHandler(http.HandlerFunc(homeView.Home)))
		se.Router.GET("/todos", apis.WrapStdHandler(http.HandlerFunc(homeView.List)))
		se.Router.POST("/todos", apis.WrapStdHandler(http.HandlerFunc(homeView.Create)))
		se.Router.PUT("/todos/{id}", apis.WrapStdHandler(http.HandlerFunc(homeView.Update)))
		se.Router.DELETE("/todos/{id}", apis.WrapStdHandler(http.HandlerFunc(homeView.Delete)))

		se.Router.GET("/api/todos/", apis.WrapStdHandler(http.HandlerFunc(todoApi.Get)))
		se.Router.POST("/api/todos/", apis.WrapStdHandler(http.HandlerFunc(todoApi.Create)))
		se.Router.PUT("/api/todos/{id}", apis.WrapStdHandler(http.HandlerFunc(todoApi.Update)))
		se.Router.DELETE("/api/todos/{id}", apis.WrapStdHandler(http.HandlerFunc(todoApi.Delete)))

		return se.Next()
	})
}
