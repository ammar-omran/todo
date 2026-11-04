package data

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

var ErrNotFound = errors.New("todo not found")

type TodoRepository struct {
	app core.App
}

func New(app core.App) *TodoRepository {
	return &TodoRepository{app: app}
}

func (r *TodoRepository) Create(ctx context.Context, description string) (Todo, error) {
	collection, err := r.app.FindCollectionByNameOrId(TodosCollection)
	if err != nil {
		return Todo{}, err
	}

	record := core.NewRecord(collection)
	record.Set("description", description)
	record.Set("done", false)

	if err := r.app.SaveWithContext(ctx, record); err != nil {
		return Todo{}, err
	}

	return fromRecord(record), nil
}

func (r *TodoRepository) GetByDone(ctx context.Context, done bool) ([]Todo, error) {
	records, err := r.app.FindRecordsByFilter(
		TodosCollection,
		"done = {:done}",
		"created",
		0,
		0,
		dbx.Params{"done": done},
	)
	if err != nil {
		return nil, err
	}

	todos := make([]Todo, len(records))
	for i, record := range records {
		todos[i] = fromRecord(record)
	}

	return todos, nil
}

func (r *TodoRepository) GetByID(ctx context.Context, id string) (Todo, error) {
	record, err := r.find(id)
	if err != nil {
		return Todo{}, err
	}

	return fromRecord(record), nil
}

func (r *TodoRepository) SetDone(ctx context.Context, id string, done bool) (Todo, error) {
	record, err := r.find(id)
	if err != nil {
		return Todo{}, err
	}

	record.Set("done", done)
	if err := r.app.SaveWithContext(ctx, record); err != nil {
		return Todo{}, err
	}

	return fromRecord(record), nil
}

func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	record, err := r.find(id)
	if err != nil {
		return err
	}

	return r.app.DeleteWithContext(ctx, record)
}

func (r *TodoRepository) find(id string) (*core.Record, error) {
	record, err := r.app.FindRecordById(TodosCollection, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return record, nil
}
