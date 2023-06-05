package task

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

var (
	Key = commons.NewKey("scope", "todo").
		Child("concept", "task")

	OperationFailedKey = Key.Child("event", "operation_failed")

	CreationRequestedKey   = Key.Child("event", "creation_requested")
	CompletionRequestedKey = Key.Child("event", "completion_requested")
	DeletionRequestedKey   = Key.Child("event", "deletion_requested")

	CreatedKey  = Key.Child("event", "created")
	DeletedKey  = Key.Child("event", "deleted")
	FinishedKey = Key.Child("event", "finished")

	OnTaskCreationRequestedKey   = Key.Child("reaktor", "on_task_creation_requested")
	OnTaskDeletionRequestedKey   = Key.Child("reaktor", "on_task_deletion_requested")
	OnTaskCompletionRequestedKey = Key.Child("reaktor", "on_task_completion_requested")

	TasksStoreKey = Key.Child("store", "tasks")
	TasksStore    = graph.NewStore(TasksStoreKey, commons.NewKey("storage", "arangodb"), "tasks")
)

func Register(env graph.Environment) error {
	err := env.RegisterConcept(graph.NewConcept(Key, graph.WithConceptDescription("A Task, something for you to do.")))
	if err != nil {
		return err
	}

	if err := env.RegisterStore(TasksStore); err != nil {
		return err
	}

	for _, event := range events() {
		if err := env.RegisterEvent(event); err != nil {
			return err
		}
	}

	for _, reaktor := range reaktors() {
		if err := env.RegisterReaktor(reaktor); err != nil {
			return err
		}
	}

	return nil
}

func events() []graph.Event {
	return []graph.Event{
		graph.NewEvent(OperationFailedKey),

		graph.NewEvent(CreationRequestedKey),
		graph.NewEvent(CompletionRequestedKey),
		graph.NewEvent(DeletionRequestedKey),

		graph.NewEvent(CreatedKey),
		graph.NewEvent(FinishedKey),
		graph.NewEvent(DeletedKey),
	}
}

func reaktors() []graph.Reaktor {
	return []graph.Reaktor{
		onTaskCreationRequestedReaktor,
		onTaskCompletionRequestedReaktor,
		onTaskDeletionRequestedReaktor,
	}
}
