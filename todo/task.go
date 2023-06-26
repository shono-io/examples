package todo

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/local"
)

func AttachTask(env *local.InventoryBuilder) {
	evtCreationRequested := inventory.NewEventReference("todo", "task", "creation_requested")
	evtCreated := inventory.NewEventReference("todo", "task", "created")
	evtCompletionRequested := inventory.NewEventReference("todo", "task", "completion_requested")
	evtFinished := inventory.NewEventReference("todo", "task", "finished")
	evtDeletionRequested := inventory.NewEventReference("todo", "task", "deletion_requested")
	evtDeleted := inventory.NewEventReference("todo", "task", "deleted")
	evtOperationFailed := inventory.NewEventReference("todo", "task", "operation_failed")

	env.
		Concept(inventory.NewConcept("todo", "task").
			Summary("A task is a single unit of work.").
			Docs(`A task is a single unit of work. It can be created, deleted and completed.`).
			Status(commons.StatusExperimental).
			Stored().
			Build()).
		Event(inventory.NewEvent("todo", "task", "operation_failed").
			Summary("Task Operation Failed").
			Docs(`An operation on a task failed`).
			Build()).
		Event(inventory.NewEvent("todo", "task", "imported").
			Summary("Task Imported").
			Docs(`A task was imported from an external source`).
			Build()).
		Event(inventory.NewEvent("todo", "task", "creation_requested").
			Summary("Task Creation Requested").
			Docs(`An external system (like an API) requested a task to be created`).
			Build()).
		Event(inventory.NewEvent("todo", "task", "created").
			Summary("Task Created").
			Docs(`A task was created`).
			Build()).
		Event(inventory.NewEvent("todo", "task", "deletion_requested").
			Summary("Task Deletion Requested").
			Docs(`An external system (like an API) requested a task to be deleted`).
			Build()).
		Event(inventory.NewEvent("todo", "task", "deleted").
			Summary("Task Deleted").
			Docs(`A task was deleted`).
			Build()).
		Reactor(inventory.NewReactor("todo", "task", "on_todo_task_creation_requested").
			Summary("Create a new task based on an external request").
			Docs(`When a creation is requested, we create the task in the store and emit a created event if it succeeds.`).
			InputEvent(evtCreationRequested).
			OutputEventCodes(evtOperationFailed.Code(), evtCreated.Code()).
			Logic(inventory.NewLogic().
				Steps(
					dsl.Log("DEBUG", "creating task ${!@} with payload ${! json(\"key\") }"),
					dsl.Transform(dsl.BloblangMapping(`
						root = this
						root.timeline.createdAt = @shono_timestamp
					`)),
					dsl.AddToStore("todo", "task", "${! json(\"key\") }"),
					dsl.Log("DEBUG", "success metadata ${!@} with payload ${! json(\"key\") }"),
					dsl.AsSuccessEvent(evtCreated, 201, `this`),
					dsl.Catch(
						dsl.Log("DEBUG", "error metadata ${!@} with payload ${! json(\"key\") }"),
						dsl.Log("ERROR", "task could not be created: ${!error()}"),
						dsl.AsFailedEvent(evtOperationFailed, 409, `"task could not be created: " + error()`),
					),
				).Test(
				inventory.NewTest("should create a task").
					When(inventory.NewTestInput(
						map[string]any{
							"key":       "1",
							"summary":   "test",
							"completed": false,
						},
						inventory.WithEventRef(evtCreationRequested),
						inventory.WithTimestamp(15))).
					Then(
						inventory.AssertMetadataContains(map[string]string{
							"shono_backbone_topic": "todo",
							"shono_status":         "201",
							"shono_kind":           "scopes/todo/concepts/task/events/created",
							"shono_timestamp":      "15",
						}),
						inventory.AssertContentEquals(map[string]interface{}{
							"key":       "1",
							"summary":   "test",
							"completed": false,
							"timeline": map[string]interface{}{
								"createdAt": "15",
							},
						})).
					Build(),
			)).
			Build()).
		Reactor(inventory.NewReactor("todo", "task", "on_todo_task_deletion_requested").
			Summary("Delete a task based on an external request").
			Docs(`When a deletion is requested, we delete the task from the store and emit a deleted event if it succeeds.`).
			InputEvent(evtDeletionRequested).
			OutputEventCodes(evtOperationFailed.Code(), evtDeleted.Code()).
			Logic(inventory.NewLogic().
				Steps(
					dsl.RemoveFromStore("todo", "task", "${! json(\"key\") }"),
					dsl.AsSuccessEvent(evtDeleted, 200, `this`),
					dsl.Catch(
						dsl.Log("ERROR", "task could not be deleted: ${!error()}"),
						dsl.AsFailedEvent(evtOperationFailed, 409, `"task could not be deleted: " + error()`),
					),
				)).
			Build()).
		Reactor(inventory.NewReactor("todo", "task", "on_todo_task_completion_requested").
			Summary("Complete a task based on an external request").
			Docs(`When a completion is requested, we complete the task in the store and emit a completed event if it succeeds.`).
			InputEvent(evtCompletionRequested).
			OutputEventCodes(evtFinished.Code(), evtOperationFailed.Code()).
			Logic(inventory.NewLogic().
				Steps(
					dsl.GetFromStore("todo", "task", "${! json(\"key\") }"),
					dsl.Transform(dsl.BloblangMapping(`
						root = this
			 			root.completed = true
						root.timeline.finishedAt = @timestamp
					`)),
					dsl.SetInStore("todo", "task", "${! json(\"key\") }"),
					dsl.AsSuccessEvent(evtFinished, 200, `this`),
					dsl.Catch(
						dsl.Log("ERROR", "task could not be completed: ${!error()}"),
						dsl.AsFailedEvent(evtOperationFailed, 409, `"task could not be completed: " + error()`),
					),
				)).
			Build())

}
