package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"to-do/data"

	"github.com/pocketbase/pocketbase/core"
)

type TodoHandler struct {
	logger *slog.Logger
	repo   *data.TodoRepository
}

type ToDoError struct {
	Code int
	Msg  string
}

func (e *ToDoError) Error() string {
	switch e.Code {
	case http.StatusBadRequest:
		e.Msg = fmt.Sprintf("Invalid Data: %s", e.Msg)
		return e.Msg
	case http.StatusUnprocessableEntity:
		e.Msg = fmt.Sprintf("Unprocessable Data: %s", e.Msg)
		return e.Msg
	default:
		return e.Msg
	}
}

func New(
	logger *slog.Logger, app core.App,
) *TodoHandler {
	return &TodoHandler{
		logger: logger,
		repo:   data.New(app),
	}
}

func (h *TodoHandler) Create(desc string) (string, error) {
	if desc == "" {
		return "", &ToDoError{Code: http.StatusBadRequest, Msg: "Description can't be empty!"}
	}

	todo, err := h.repo.Create(context.Background(), desc)
	if err != nil {
		return "", mapError(err)
	}

	return todo.Id, nil
}

func (h *TodoHandler) Update(status string, id string) error {
	if id == "" {
		return &ToDoError{Code: http.StatusBadRequest, Msg: "Invalid Id"}
	}

	completed := status != "reopen"
	_, err := h.repo.SetDone(context.Background(), id, completed)

	return mapError(err)
}

type ToDoResponse struct {
	Id          string    `json:"ID"`
	Description string    `json:"Description"`
	Created     time.Time `json:"CreatedAt"`
	Done        bool      `json:"Done"`
}

func (h *TodoHandler) Get(status string) ([]ToDoResponse, error) {
	completed := status == "completed"
	res, err := h.repo.GetByDone(context.Background(), completed)
	if err != nil {
		return nil, mapError(err)
	}

	return Map(res, func(td data.Todo) ToDoResponse {
		return ToDoResponse{td.Id, td.Description, td.Created, td.Done}
	}), nil
}

func (h *TodoHandler) GetByID(id string) (ToDoResponse, error) {
	td, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		return ToDoResponse{}, mapError(err)
	}

	return ToDoResponse{td.Id, td.Description, td.Created, td.Done}, nil
}

func (h *TodoHandler) Delete(id string) error {
	if id == "" {
		return &ToDoError{Code: http.StatusBadRequest, Msg: "Invalid Id"}
	}

	return mapError(h.repo.Delete(context.Background(), id))
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, data.ErrNotFound) {
		return &ToDoError{Code: http.StatusNotFound, Msg: "Todo not found"}
	}

	return err
}

func Map[T any, R any](items []T, fn func(T) R) []R {
	result := make([]R, len(items))

	for i, item := range items {
		result[i] = fn(item)
	}

	return result
}
