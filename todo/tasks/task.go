package tasks

import (
	"github.com/shono-io/shono/graph"
)

var scope = graph.NewScopeSpec("todo").
	Experimental().
	Summary("Manage todos at scale").
	Concept(graph.NewConceptSpec("task").
		Summary("A Task, something for you to do.").
		Experimental())

func Reaktors() ([]graph.Reaktor, error) {
	var result []graph.Reaktor

	r, err := graph.
		InputEvent("todo", "task", "creation_requested").
		ExecuteFor("todo", "task", onTaskCreationRequestedLogic()...).
		Producing("operation_failed", "when the system was unable to create the task.").
		Producing("created", "when the system was able to create the task.").
		NamedAs("OnTaskCreationRequested").
		Build()
	if err != nil {
		return nil, err
	}
	result = append(result, *r)

	r, err = graph.
		InputEvent("todo", "task", "deletion_requested").
		ExecuteFor("todo", "task", onTaskDeletionRequestedLogic()...).
		Producing("operation_failed", "when the system was unable to remove the task.").
		Producing("deleted", "when the system has removed the task.").
		NamedAs("OnTaskDeletionRequested").
		Build()
	if err != nil {
		return nil, err
	}
	result = append(result, *r)

	r, err = graph.
		InputEvent("todo", "task", "completion_requested").
		ExecuteFor("todo", "task", onTaskCompletionRequestedLogic()...).
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
