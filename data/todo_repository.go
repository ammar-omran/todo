package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const createToDo = `-- name: CreateToDo :one
INSERT INTO ToDos (description, created, done)
VALUES ($1, $2, $3)
RETURNING id, description, created, done`

type CreateToDoParams struct {
	Description string
	Created     time.Time
	Done        bool
}

func (q *Queries) CreateToDo(ctx context.Context, arg CreateToDoParams) (ToDo, error) {
	row := q.db.QueryRowContext(ctx, createToDo,
		arg.Description,
		arg.Created,
		arg.Done)
	var i ToDo
	err := row.Scan(
		&i.Id,
		&i.Description,
		&i.Created,
		&i.Done,
	)
	return i, err
}

const changeToDoStatus = `-- name: changeToDoStatus :one
	UPDATE Todos
	SET done = $2
	WHERE id = $1
	RETURNING id;
`

func (q *Queries) ChangeToDoStatus(ctx context.Context, id int, status bool) error {
	var updatedID int

	err := q.db.QueryRowContext(
		ctx,
		changeToDoStatus,
		id,
		status,
	).Scan(&updatedID)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return err
}

const getToDosByState = `-- name: getToDosByState :many
SELECT Id, Description, Created, Done
	FROM ToDos
  WHERE Done = $1
`

func (q *Queries) GetToDosByState(ctx context.Context, done bool) ([]ToDo, error) {
	rows, err := q.db.QueryContext(ctx, getToDosByState, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []ToDo

	for rows.Next() {
		var todo ToDo

		err := rows.Scan(
			&todo.Id,
			&todo.Description,
			&todo.Created,
			&todo.Done,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

const deleteTodo = `-- name: DeleteTodo :one
	DELETE FROM Todos
	WHERE id = $1
	RETURNING id;
`

func (q *Queries) DeleteTodo(ctx context.Context, id int) error {
	var updatedID int

	err := q.db.QueryRowContext(
		ctx,
		deleteTodo,
		id,
	).Scan(&updatedID)

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return err
}
