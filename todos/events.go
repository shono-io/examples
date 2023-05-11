package todos

import (
	go_shono "github.com/shono-io/go-shono"
)

var (
	org     = "my-org"
	space   = "todos"
	concept = "todo"

	TaskAdded    = go_shono.NewEvent(go_shono.NewEventId(org, space, concept, "added"), new(TaskAddedEvent))
	TaskFinished = go_shono.NewEvent(go_shono.NewEventId(org, space, concept, "finished"), new(TaskFinishedEvent))
	TaskDeleted  = go_shono.NewEvent(go_shono.NewEventId(org, space, concept, "deleted"), new(TaskDeletedEvent))
)

type TaskAddedEvent struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type TaskFinishedEvent struct {
	Id string `json:"id"`
}

type TaskDeletedEvent struct {
	Id string `json:"id"`
}
