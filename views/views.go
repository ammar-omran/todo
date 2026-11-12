package views

import (
	"io/fs"
	"net/http"
	"time"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/template"
)

type todoView struct {
	Id          string
	Description string
	Created     time.Time
	Done        bool
}

type errorMessage struct {
	StatusMessage string
}

// Views holds the page and fragment templates used by the HTTP handlers.
type Views struct {
	fsys     fs.FS
	registry *template.Registry
}

// New creates a Views instance backed by the provided template filesystem.
func New(fsys fs.FS) *Views {
	return &Views{
		fsys:     fsys,
		registry: template.NewRegistry(),
	}
}

// Home renders the todo page.
func (v *Views) Home(e *core.RequestEvent) error {
	html, err := v.registry.LoadFS(v.fsys, "templates/home/home.html").Render(nil)
	if err != nil {
		return e.InternalServerError("Failed to render home page.", err)
	}

	return e.HTML(http.StatusOK, html)
}

func (v *Views) todoListHTML(todos []todoView) (string, error) {
	return v.registry.LoadFS(
		v.fsys,
		"templates/home/partials/todo-list.html",
		"templates/home/partials/todo-row.html",
	).Render(todos)
}

func (v *Views) errorMessageHTML(message string) (string, error) {
	return v.registry.LoadFS(
		v.fsys,
		"templates/shared/partials/error-message.html",
	).Render(errorMessage{StatusMessage: message})
}
