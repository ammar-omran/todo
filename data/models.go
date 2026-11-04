package data

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
)

const TodosCollection = "todos"

type Todo struct {
	Id          string
	Description string
	Created     time.Time
	Done        bool
}

func fromRecord(record *core.Record) Todo {
	return Todo{
		Id:          record.Id,
		Description: record.GetString("description"),
		Created:     record.GetDateTime("created").Time(),
		Done:        record.GetBool("done"),
	}
}
