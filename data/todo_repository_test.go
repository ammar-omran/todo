package data

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE ToDos (
			Id INTEGER PRIMARY KEY AUTOINCREMENT,
			Description TEXT NOT NULL,
			Created DATETIME NOT NULL,
			Done BOOLEAN NOT NULL
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func setupQueries(t *testing.T) *Queries {
	t.Helper()

	db := setupTestDB(t)

	return New(db)
}

func TestCreateToDo(t *testing.T) {
	testCases := []struct {
		desc        string
		description string
		done        bool
	}{
		{
			desc:        "create incomplete todo",
			description: "Learn Go testing",
			done:        false,
		},
		{
			desc:        "create completed todo",
			description: "Learn SQL",
			done:        true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			q := setupQueries(t)

			created := time.Now()

			todo, err := q.CreateToDo(context.Background(), CreateToDoParams{
				Description: tC.description,
				Created:     created,
				Done:        tC.done,
			})

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if todo.Id == 0 {
				t.Error("expected generated id")
			}

			if todo.Description != tC.description {
				t.Errorf(
					"expected description %q, got %q",
					tC.description,
					todo.Description,
				)
			}

			if todo.Done != tC.done {
				t.Errorf(
					"expected done %v, got %v",
					tC.done,
					todo.Done,
				)
			}
		})
	}
}

func TestChangeToDoStatus(t *testing.T) {
	testCases := []struct {
		desc        string
		id          int
		status      bool
		expectError bool
	}{
		{
			desc:        "change existing todo",
			status:      true,
			expectError: false,
		},
		{
			desc:        "todo does not exist",
			id:          -1,
			status:      true,
			expectError: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			q := setupQueries(t)

			todo, err := q.CreateToDo(context.Background(), CreateToDoParams{
				Description: "Learn TDD",
				Created:     time.Now(),
				Done:        false,
			})
			if err != nil {
				t.Fatal(err)
			}

			id := todo.Id

			if tC.expectError {
				id = tC.id
			}

			err = q.ChangeToDoStatus(
				context.Background(),
				id,
				tC.status,
			)

			if tC.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetToDosByState(t *testing.T) {
	q := setupQueries(t)

	_, err := q.CreateToDo(context.Background(), CreateToDoParams{
		Description: "Incomplete task",
		Created:     time.Now(),
		Done:        false,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = q.CreateToDo(context.Background(), CreateToDoParams{
		Description: "Completed task",
		Created:     time.Now(),
		Done:        true,
	})
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc          string
		done          bool
		expectedCount int
	}{
		{
			desc:          "get incomplete todos",
			done:          false,
			expectedCount: 1,
		},
		{
			desc:          "get completed todos",
			done:          true,
			expectedCount: 1,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			todos, err := q.GetToDosByState(
				context.Background(),
				tC.done,
			)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(todos) != tC.expectedCount {
				t.Errorf(
					"expected %d todos, got %d",
					tC.expectedCount,
					len(todos),
				)
			}

			for _, todo := range todos {
				if todo.Done != tC.done {
					t.Errorf(
						"expected done=%v, got %v",
						tC.done,
						todo.Done,
					)
				}
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	testCases := []struct {
		desc        string
		expectError bool
	}{
		{
			desc:        "delete existing todo",
			expectError: false,
		},
		{
			desc:        "delete nonexistent todo",
			expectError: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			q := setupQueries(t)

			todo, err := q.CreateToDo(context.Background(), CreateToDoParams{
				Description: "Delete me",
				Created:     time.Now(),
				Done:        false,
			})

			if err != nil {
				t.Fatal(err)
			}

			id := todo.Id

			if tC.expectError {
				id = 999
			}

			// Act
			err = q.DeleteTodo(context.Background(), id)

			// Assert
			if tC.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify it was actually deleted.
			_, err = q.GetToDosByState(
				context.Background(),
				false,
			)

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
