package data

import (
	"time"
)

type ToDo struct {
	Id          int
	Description string
	Created     time.Time
	Done        bool
}

func CreateToDo(description string) ToDo {
	return ToDo{
		Id:          0,
		Description: description,
		Created:     time.Now(),
		Done:        false,
	}
}
