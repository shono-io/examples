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
					dsl.AddToStore("todo", "task", "${! json(\"key\") }", dsl.BloblangMapping(`
						root.key = this.key
						root.summary = this.summary
						root.completed = this.completed
						root.timeline.createdAt = @kafka_timestamp_unix
					`)),
					dsl.Transform(dsl.BloblangMapping(`
						meta io_shono_kind = "scopes/todo/concepts/task/events/created"
						root.status = 201
						root.task = this
					`)),
					dsl.Catch(
						dsl.Log("ERROR", "task could not be created: ${!error()}"),
						dsl.Transform(dsl.BloblangMapping(`
							meta io_shono_kind = "scopes/todo/concepts/task/events/operation_failed"
							root.status = 409
							root.message = "task could not be created: ${!error()}"
						`)),
					),
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
					dsl.Transform(dsl.BloblangMapping(`
						meta io_shono_kind = "scopes/todo/concepts/task/events/deleted"
						root.status = 200
						root.removed = this
					`)),
					dsl.Catch(
						dsl.Log("ERROR", "task could not be deleted: ${!error()}"),
						dsl.Transform(dsl.BloblangMapping(`
							meta io_shono_kind = "scopes/todo/concepts/task/events/operation_failed"
							root.status = 409
							root.message = "task could not be deleted: ${!error()}"
						`)),
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
					dsl.SetInStore("todo", "task", "${! json(\"key\") }", dsl.BloblangMapping(`
						root = this
						root.finished = true
					`)),
					dsl.Transform(dsl.BloblangMapping(`
						meta io_shono_kind = "scopes/todo/concepts/task/events/finished"
						root = this
					`)),
					dsl.Catch(
						dsl.Log("ERROR", "task could not be completed: ${!error()}"),
						dsl.Transform(dsl.BloblangMapping(`
							meta io_shono_kind = "scopes/todo/concepts/task/events/operation_failed"
							root.status = 409
							root.message = "task could not be finished: ${!error()}"
						`)),
					),
				)).
			Build())

}
