package todos

import (
	go_shono "github.com/shono-io/go-shono"
)

var (
	org     = "my-org"
	domain  = "todos"
	concept = "todo"

	AddTodo      = go_shono.NewEvent(org, domain, concept, "add", new(AddTodoEvent))
	TodoAdded    = go_shono.NewEvent(org, domain, concept, "added", new(TodoAddedEvent))
	FinishTodo   = go_shono.NewEvent(org, domain, concept, "finish", new(FinishTodoEvent))
	TodoFinished = go_shono.NewEvent(org, domain, concept, "finished", new(TodoFinishedEvent))
	DeleteTodo   = go_shono.NewEvent(org, domain, concept, "delete", new(DeleteTodoEvent))
	TodoDeleted  = go_shono.NewEvent(org, domain, concept, "deleted", new(TodoDeletedEvent))
)

type AddTodoEvent struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type TodoAddedEvent struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

type FinishTodoEvent struct {
	Id string `json:"id"`
}

type TodoFinishedEvent struct {
	Id string `json:"id"`
}

type DeleteTodoEvent struct {
	Id string `json:"id"`
}

type TodoDeletedEvent struct {
	Id string `json:"id"`
}
