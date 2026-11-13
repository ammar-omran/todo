package views

import (
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"to-do/internal/i18n"
	"to-do/internal/render"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type todoView struct {
	Id          string
	Description string
	Created     time.Time
	Done        bool
}

type pageView struct {
	Lang string
	Dir  string
}

type todoRowView struct {
	Todo         todoView
	CreatedLabel string
}

type todoListView struct {
	Rows []todoRowView
}

type errorMessage struct {
	StatusMessage string
}

// Views holds the page and fragment templates used by the HTTP handlers.
type Views struct {
	fsys     fs.FS
	registry *render.Registry
}

// New creates a Views instance backed by the provided template filesystem.
func New(fsys fs.FS) *Views {
	return &Views{
		fsys:     fsys,
		registry: render.NewRegistry(),
	}
}

// Home renders the todo page.
func (v *Views) Home(e *core.RequestEvent) error {
	tr := translatorFor(e)

	html, err := v.registry.
		LoadFS(v.fsys, "templates/index.html").
		Render(pageView{Lang: tr.Lang, Dir: tr.Dir}, funcsFor(tr))
	if err != nil {
		return e.InternalServerError("Failed to render home page.", err)
	}

	return e.HTML(http.StatusOK, html)
}

func (v *Views) todoListHTML(rows []todoRowView, tr i18n.Translator) (string, error) {
	return v.registry.
		LoadFS(
			v.fsys,
			"templates/home/todo-list.html",
			"templates/home/todo-row.html",
		).
		Render(todoListView{Rows: rows}, funcsFor(tr))
}

func (v *Views) todoRowHTML(todo todoView, tr i18n.Translator) (string, error) {
	return v.registry.
		LoadFS(v.fsys, "templates/home/todo-row.html").
		Render(todoRowView{
			Todo:         todo,
			CreatedLabel: tr.Date(todo.Created),
		}, funcsFor(tr))
}

func (v *Views) errorMessageHTML(message string, tr i18n.Translator) (string, error) {
	return v.registry.
		LoadFS(v.fsys, "templates/shared/partials/error-message.html").
		Render(errorMessage{StatusMessage: message}, funcsFor(tr))
}

func funcsFor(tr i18n.Translator) template.FuncMap {
	return template.FuncMap{
		"T": tr.T,
	}
}
