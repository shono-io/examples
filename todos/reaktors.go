package todos

import (
	context2 "context"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/resources"
	"github.com/shono-io/shono-examples/arangodb"
)

func Reaktors(runtime *go_shono.Runtime) []go_shono.Reaktor {
	return []go_shono.Reaktor{
		go_shono.MustNewReaktor("addTask", go_shono.ListenFor(TaskAdded), go_shono.WithHandler(addTodo(runtime))),
		go_shono.MustNewReaktor("finishTask", go_shono.ListenFor(TaskFinished), go_shono.WithHandler(finishTodo(runtime))),
		go_shono.MustNewReaktor("deleteTask", go_shono.ListenFor(TaskDeleted), go_shono.WithHandler(deleteTodo(runtime))),
	}
}

func addTodo(runtime *go_shono.Runtime) go_shono.ReaktorFunc {
	return func(ctx context2.Context, value any, w go_shono.Writer) {
		evt := value.(*TaskAddedEvent)

		db := runtime.Resource("db").(resources.DocumentStore[arangodb.Query])
		db.MustSet(ctx, "tasks", evt.Id, Task{
			Id:     evt.Id,
			Label:  evt.Label,
			IsDone: false,
		})
	}
}

func finishTodo(runtime *go_shono.Runtime) go_shono.ReaktorFunc {
	return func(ctx context2.Context, value any, w go_shono.Writer) {
		evt := value.(*TaskFinishedEvent)
		db := runtime.Resource("db").(resources.DocumentStore[arangodb.Query])

		// -- get the task
		var task Task
		if fnd := db.MustGet(ctx, "tasks", evt.Id, &task); !fnd {
			panic("task with id " + evt.Id + " not found")
		}

		// -- update the task
		task.IsDone = true
		db.MustSet(ctx, "tasks", evt.Id, task)
	}
}

func deleteTodo(runtime *go_shono.Runtime) go_shono.ReaktorFunc {
	return func(ctx context2.Context, value any, w go_shono.Writer) {
		evt := value.(*TaskDeletedEvent)
		db := runtime.Resource("db").(resources.DocumentStore[arangodb.Query])

		// -- remove the task
		db.MustDelete(ctx, "tasks", evt.Id)
	}
}
