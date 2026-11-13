package app

import (
	"fmt"
	"io/fs"
	"os"
	"to-do/views"

	"ganigeorgiev.com/pocketbase/pocketbase"
	"ganigeorgiev.com/pocketbase/pocketbase/apis"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/plugins/migratecmd"
)

type App struct {
	pb   *pocketbase.PocketBase
	view *views.Views
}

func New(templates fs.FS) *App {
	pb := pocketbase.NewWithConfig(pocketbase.Config{
		HideStartBanner: false,
		DefaultDataDir:  "./pb_data",
	})

	migratecmd.MustRegister(pb, pb.RootCmd, migratecmd.Config{
		Automigrate:  true,
		Dir:          "./migrations",
		TemplateLang: migratecmd.TemplateLangGo,
	})

	return &App{
		pb:   pb,
		view: views.New(templates),
	}
}

func (a *App) Start() error {
	a.pb.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/", a.view.Home)
		se.Router.POST("/language", a.view.SetLanguage)
		se.Router.GET("/todos", a.view.List)
		se.Router.POST("/todos", a.view.Create)
		se.Router.PATCH("/todos/{id}", a.view.Toggle)
		se.Router.DELETE("/todos/{id}", a.view.Delete)

		se.Router.GET("/static/{path...}", apis.Static(os.DirFS("./static"), false))
		return se.Next()
	})

	if len(os.Args) == 1 {
		a.pb.RootCmd.SetArgs([]string{"serve", "--http=localhost:8080"})
	}

	if err := a.pb.Start(); err != nil {
		return fmt.Errorf("start pocketbase: %w", err)
	}

	return nil
}
