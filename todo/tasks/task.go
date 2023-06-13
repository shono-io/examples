package tasks

import (
	"github.com/shono-io/shono/graph"
)

func Reaktors() ([]graph.Reaktor, error) {
	var result []graph.Reaktor

	r, err := graph.
		InputEvent("todo", "task", "creation_requested").
		ExecuteFor("todo", "tasks", onTaskCreationRequestedLogic()...).
		Producing("operation_failed", "when the system was unable to create the task.").
		Producing("created", "when the system was able to create the task.").
		NamedAs("OnTaskCreationRequested").
		Build()
	if err != nil {
		return nil, err
	}
	result = append(result, *r)

	r, err = graph.
		InputEvent("todo", "tasks", "deletion_requested").
		ExecuteFor("todo", "tasks", onTaskDeletionRequestedLogic()...).
		Producing("operation_failed", "when the system was unable to remove the task.").
		Producing("deleted", "when the system has removed the task.").
		NamedAs("OnTaskDeletionRequested").
		Build()
	if err != nil {
		return nil, err
	}
	result = append(result, *r)

	r, err = graph.
		InputEvent("todo", "tasks", "completion_requested").
		ExecuteFor("todo", "tasks", onTaskCompletionRequestedLogic()...).
		Producing("operation_failed", "when the system was unable to complete the task.").
		Producing("finished", "when the task has been finished.").
		NamedAs("OnTaskCompletionRequested").
		Build()
	if err != nil {
		return nil, err
	}
	result = append(result, *r)

	return result, nil
}
