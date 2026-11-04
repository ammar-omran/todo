package views

import (
	"errors"
	"html/template"
	"net/http"
	"strings"
	"to-do/handler"
)

type Home struct {
	Th   *handler.TodoHandler
	Tmpl *template.Template
}

type indexPage struct {
	Todos []handler.ToDoResponse
}

type errorPage struct {
	StatusCode    int
	StatusMessage string
}

type errorMessage struct {
	StatusMessage string
}

func (a *Home) Home(w http.ResponseWriter, req *http.Request) {
	var status = req.URL.Query().Get("status")
	data, err := a.Th.Get(status)
	if err != nil {
		a.renderError(w, req, err)
		return
	}

	a.render(w, http.StatusOK, "index.html", indexPage{
		Todos: data,
	})
}

func (a *Home) List(w http.ResponseWriter, req *http.Request) {
	var status = req.URL.Query().Get("status")
	data, err := a.Th.Get(status)
	if err != nil {
		a.renderError(w, req, err)
		return
	}

	a.render(w, http.StatusOK, "todo-list", data)
}

func (a *Home) Create(w http.ResponseWriter, req *http.Request) {
	desc := strings.TrimSpace(req.FormValue("description"))
	id, err := a.Th.Create(desc)
	if err != nil {
		a.renderError(w, req, err)
		return
	}

	todo, err := a.Th.GetByID(id)
	if err != nil {
		a.renderError(w, req, err)
		return
	}

	a.render(w, http.StatusCreated, "todo-row", todo)
}

func (a *Home) Update(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	status := req.URL.Query().Get("status")

	if err := a.Th.Update(status, id); err != nil {
		a.renderError(w, req, err)
		return
	}

	todo, err := a.Th.GetByID(id)
	if err != nil {
		a.renderError(w, req, err)
		return
	}

	a.render(w, http.StatusOK, "todo-row", todo)
}

func (a *Home) Delete(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	if err := a.Th.Delete(id); err != nil {
		a.renderError(w, req, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (a *Home) render(w http.ResponseWriter, code int, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	_ = a.Tmpl.ExecuteTemplate(w, name, data)
}

func (a *Home) renderError(w http.ResponseWriter, req *http.Request, err error) {
	code := http.StatusInternalServerError
	msg := "Something went wrong"

	if te, ok := errors.AsType[*handler.ToDoError](err); ok {
		code = te.Code
		msg = te.Error()
	}

	if req.Header.Get("HX-Request") == "true" {
		a.render(w, code, "error-message", errorMessage{StatusMessage: msg})
		return
	}

	a.render(w, code, "error.html", errorPage{
		StatusCode:    code,
		StatusMessage: msg,
	})
}
