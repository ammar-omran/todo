package app

import (
	"html/template"
	"net/http"
	"os"
	"to-do/api"
	"to-do/handler"
	"to-do/views"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func (a *App) loadRoutes(tmpl *template.Template) {
	todoHandler := handler.New(a.logger, a.pb)
	todoApi := &api.Api{Th: todoHandler}
	homeView := &views.Home{Th: todoHandler, Tmpl: tmpl}

	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("GET /", todoApi.Get)
	apiRouter.HandleFunc("POST /", todoApi.Create)
	apiRouter.HandleFunc("PUT /{id}", todoApi.Update)
	apiRouter.HandleFunc("DELETE /{id}", todoApi.Delete)

	a.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.InstallerFunc = nil

		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./static"), false))

		se.Router.GET("/", apis.WrapStdHandler(http.HandlerFunc(homeView.Home)))
		se.Router.GET("/todos", apis.WrapStdHandler(http.HandlerFunc(homeView.List)))
		se.Router.POST("/todos", apis.WrapStdHandler(http.HandlerFunc(homeView.Create)))
		se.Router.PUT("/todos/{id}", apis.WrapStdHandler(http.HandlerFunc(homeView.Update)))
		se.Router.DELETE("/todos/{id}", apis.WrapStdHandler(http.HandlerFunc(homeView.Delete)))

		apiHandler := apis.WrapStdHandler(http.StripPrefix("/api/todos", apiRouter))
		se.Router.GET("/api/todos/{path...}", apiHandler)
		se.Router.POST("/api/todos/{path...}", apiHandler)
		se.Router.PUT("/api/todos/{path...}", apiHandler)
		se.Router.DELETE("/api/todos/{path...}", apiHandler)

		return se.Next()
	})
}
