package views

import (
	"net/http"
	"strings"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/forms"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/router"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/search"
)

const (
	todosCollectionName = "todos"
	hxTriggerHeader     = "HX-Trigger"
)

// List renders the todo list fragment.
func (v *Views) List(e *core.RequestEvent) error {
	collection, err := v.todosCollection(e)
	if err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	requestInfo, err := e.RequestInfo()
	if err != nil {
		return v.renderError(e, router.NewBadRequestError("Failed to read request.", err))
	}

	if !requestInfo.HasSuperuserAuth() && collection.ListRule == nil {
		return v.renderError(e, e.ForbiddenError("Only superusers can perform this action.", nil))
	}

	query := e.App.RecordQuery(collection)

	if !requestInfo.HasSuperuserAuth() && collection.ListRule != nil && *collection.ListRule != "" {
		resolver := core.NewRecordFieldResolver(e.App, collection, requestInfo, true)

		expr, err := search.FilterData(*collection.ListRule).BuildExpr(resolver)
		if err != nil {
			return v.renderError(e, router.NewBadRequestError("Failed to build list rule.", err))
		}
		query.AndWhere(expr)

		if err := resolver.UpdateQuery(query); err != nil {
			return v.renderError(e, router.NewBadRequestError("Failed to apply list rule.", err))
		}
	}

	switch e.Request.URL.Query().Get("filter") {
	case "active", "(done=false)":
		query.AndWhere(dbx.HashExp{"done": false})
	case "completed", "(done=true)":
		query.AndWhere(dbx.HashExp{"done": true})
	}

	records := []*core.Record{}
	if err := query.OrderBy("created ASC").Limit(200).All(&records); err != nil {
		return v.renderError(e, router.NewInternalServerError("Failed to load todos.", err))
	}

	html, err := v.todoListHTML(todoViews(records))
	if err != nil {
		return v.renderError(e, router.NewInternalServerError("Failed to render todos.", err))
	}

	return e.HTML(http.StatusOK, html)
}

// Create adds a new todo.
func (v *Views) Create(e *core.RequestEvent) error {
	collection, err := v.todosCollection(e)
	if err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	requestInfo, err := e.RequestInfo()
	if err != nil {
		return v.renderError(e, router.NewBadRequestError("Failed to read request.", err))
	}

	description, _ := requestInfo.Body["description"].(string)

	record := core.NewRecord(collection)
	form := forms.NewRecordUpsert(e.App, record)
	form.Load(map[string]any{"description": strings.TrimSpace(description)})

	allowed, err := canCreate(e, collection, record, requestInfo)
	if err != nil {
		return v.renderError(e, router.ToApiError(err))
	}
	if !allowed {
		return v.renderError(e, e.ForbiddenError("You are not allowed to perform this request.", nil))
	}

	if err := form.Submit(); err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	e.Response.Header().Set(hxTriggerHeader, "todos-changed")

	return e.NoContent(http.StatusNoContent)
}

// Toggle flips the done state of a todo.
func (v *Views) Toggle(e *core.RequestEvent) error {
	collection, err := v.todosCollection(e)
	if err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	record, err := e.App.FindRecordById(collection, e.Request.PathValue("id"))
	if err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	requestInfo, err := e.RequestInfo()
	if err != nil {
		return v.renderError(e, router.NewBadRequestError("Failed to read request.", err))
	}

	if !requestInfo.HasSuperuserAuth() {
		allowed, err := e.App.CanAccessRecord(record, requestInfo, collection.UpdateRule)
		if err != nil {
			return v.renderError(e, router.ToApiError(err))
		}
		if !allowed {
			return v.renderError(e, e.ForbiddenError("You are not allowed to perform this request.", nil))
		}
	}

	record.Set("done", !record.GetBool("done"))

	if err := e.App.Save(record); err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	e.Response.Header().Set(hxTriggerHeader, "todos-changed")

	return e.NoContent(http.StatusNoContent)
}

// Delete removes a todo.
func (v *Views) Delete(e *core.RequestEvent) error {
	collection, err := v.todosCollection(e)
	if err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	record, err := e.App.FindRecordById(collection, e.Request.PathValue("id"))
	if err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	requestInfo, err := e.RequestInfo()
	if err != nil {
		return v.renderError(e, router.NewBadRequestError("Failed to read request.", err))
	}

	if !requestInfo.HasSuperuserAuth() {
		allowed, err := e.App.CanAccessRecord(record, requestInfo, collection.DeleteRule)
		if err != nil {
			return v.renderError(e, router.ToApiError(err))
		}
		if !allowed {
			return v.renderError(e, e.ForbiddenError("You are not allowed to perform this request.", nil))
		}
	}

	if err := e.App.Delete(record); err != nil {
		return v.renderError(e, router.ToApiError(err))
	}

	e.Response.Header().Set(hxTriggerHeader, "todos-changed")

	return e.NoContent(http.StatusNoContent)
}

func (v *Views) todosCollection(e *core.RequestEvent) (*core.Collection, error) {
	collection, err := e.App.FindCachedCollectionByNameOrId(todosCollectionName)
	if err != nil {
		return nil, e.NotFoundError("Missing todos collection.", err)
	}

	return collection, nil
}

func canCreate(e *core.RequestEvent, collection *core.Collection, record *core.Record, requestInfo *core.RequestInfo) (bool, error) {
	if requestInfo.HasSuperuserAuth() {
		return true, nil
	}

	if collection.CreateRule == nil {
		return false, nil
	}

	if *collection.CreateRule == "" {
		return true, nil
	}

	var allowed bool

	err := forms.NewRecordUpsert(e.App, record).DrySubmit(func(txApp core.App, dryRecord *core.Record) error {
		ok, err := txApp.CanAccessRecord(dryRecord, requestInfo, collection.CreateRule)
		allowed = ok
		return err
	})

	return allowed, err
}

func todoViews(records []*core.Record) []todoView {
	todos := make([]todoView, 0, len(records))

	for _, record := range records {
		todos = append(todos, todoView{
			Id:          record.Id,
			Description: record.GetString("description"),
			Created:     record.GetDateTime("created").Time(),
			Done:        record.GetBool("done"),
		})
	}

	return todos
}

func (v *Views) renderError(e *core.RequestEvent, apiErr *router.ApiError) error {
	html, err := v.errorMessageHTML(errorText(apiErr))
	if err != nil {
		return err
	}

	e.Response.Header().Set("HX-Retarget", "#form-error")
	e.Response.Header().Set("HX-Reswap", "innerHTML")

	return e.HTML(apiErr.Status, html)
}

func errorText(apiErr *router.ApiError) string {
	for _, value := range apiErr.Data {
		field, ok := value.(map[string]any)
		if !ok {
			continue
		}

		if message, ok := field["message"].(string); ok && message != "" {
			return message
		}
	}

	return apiErr.Message
}
