package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"to-do/data"
)

type TodoHandler struct {
	logger *slog.Logger
	repo   *data.Queries
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
	logger *slog.Logger, db *sql.DB,
) *TodoHandler {
	return &TodoHandler{
		logger: logger,
		repo:   data.New(db),
	}
}

func (h *TodoHandler) Create(desc string) (int, error) {
	if desc == "" {
		return -1, &ToDoError{Code: http.StatusBadRequest, Msg: "Description can't be empty!"}
	}
	d := data.CreateToDo(desc)
	res, err := h.repo.CreateToDo(
		context.Background(), data.CreateToDoParams{
			Description: d.Description,
			Created:     d.Created,
			Done:        d.Done,
		},
	)
	if err != nil {
		return -1, err
	}

	return res.Id, nil
}

func (h *TodoHandler) Update(status string, id string) error {
	completed := status != "reopen"
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &ToDoError{Code: http.StatusUnprocessableEntity, Msg: "Invalid Id value, only numpers Allowed"}
	}
	if err := h.repo.ChangeToDoStatus(context.Background(), idInt, completed); err != nil {
		return err
	}
	return nil
}

type ToDoResponse struct {
	Id          int       `json:"ID"`
	Description string    `json:"Description"`
	Created     time.Time `json:"CreatedAt"`
	Done        bool      `json:"Done"`
}

func (h *TodoHandler) Get(status string) ([]ToDoResponse, error) {
	completed := status == "completed"
	res, err := h.repo.GetToDosByState(context.Background(), completed)
	if err != nil {
		return nil, err
	}
	return Map(res, func(td data.ToDo) ToDoResponse {
		return ToDoResponse{td.Id, td.Description, td.Created, td.Done}
	}), nil
}

func (h *TodoHandler) Delete(id string) error {
	if id == "" {
		return &ToDoError{Code: http.StatusBadRequest, Msg: "Invalid Id"}
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &ToDoError{Code: http.StatusUnprocessableEntity, Msg: "Invalid Id value, only numpers Allowed"}
	}
	err = h.repo.DeleteTodo(context.Background(), int(idInt))
	if err != nil {
		return err
	}
	return nil
}

func Map[T any, R any](items []T, fn func(T) R) []R {
	result := make([]R, len(items))

	for i, item := range items {
		result[i] = fn(item)
	}

	return result
}
