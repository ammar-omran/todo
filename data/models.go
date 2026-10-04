package data

import (
	"time"
)

type ToDo struct {
	Id          int       `json:"ID"`
	Description string    `json:"Description"`
	Created     time.Time `json:"CreatedAt"`
	Done        bool      `json:"Done"`
}

func CreateToDo(description string) ToDo {
	return ToDo{
		Id:          0,
		Description: description,
		Created:     time.Now(),
		Done:        false,
	}
}
