package todos

import (
	"context"
	"errors"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/resources"
	"github.com/shono-io/shono-examples/arangodb"
)

func Reaktors(runtime *go_shono.Runtime) []go_shono.Reaktor {
	return []go_shono.Reaktor{
		go_shono.MustNewReaktor("addTodo", go_shono.ListenFor(AddTodo), go_shono.WithHandler(addTodo(runtime))),
		go_shono.MustNewReaktor("finishTodo", go_shono.ListenFor(FinishTodo), go_shono.WithHandler(finishTodo(runtime))),
		go_shono.MustNewReaktor("deleteTodo", go_shono.ListenFor(DeleteTodo), go_shono.WithHandler(deleteTodo(runtime))),
	}
}

func addTodo(runtime *go_shono.Runtime) go_shono.ReaktorFunc {
	return func(ctx context.Context, value any, w go_shono.Writer) {
		cmd := value.(AddTodoEvent)

		db := runtime.Resource("db").(resources.DocumentStore[arangodb.Query])
		db.MustSet(ctx, "todos", cmd.Id, Todo{
			Id:     cmd.Id,
			Label:  cmd.Label,
			IsDone: false,
		})

		w.Write(TodoAdded, go_shono.KeyFromContext(ctx), TodoAddedEvent{
			Id: cmd.Id,
		})
	}
}

func finishTodo(runtime *go_shono.Runtime) go_shono.ReaktorFunc {
	return func(ctx context.Context, value any, w go_shono.Writer) {
		cmd := evt.(FinishTodoEvent)
		db := runtime.Resource("db").(resources.DocumentStore[arangodb.Query])

		// -- get the todo
		var todo Todo
		if fnd := db.MustGet(ctx, "todos", cmd.Id, &todo); !fnd {
			ctx.Failed(TodoFinished, errors.New("todo not found"))
			return
		}

		// -- update the todo
		todo.IsDone = true
		db.MustSet(ctx, "todos", cmd.Id, todo)

		// -- send the event
		ctx.Send(TodoFinished, TodoFinishedEvent{
			Id: cmd.Id,
		})
	}
}

func deleteTodo(runtime *go_shono.Runtime) go_shono.ReaktorFunc {
	return func(ctx context.Context, value any, w go_shono.Writer) {
		cmd := evt.(DeleteTodoEvent)
		db := runtime.Resource("db").(resources.DocumentStore[arangodb.Query])

		// -- remove the todo
		db.MustDelete(ctx, "todos", cmd.Id)

		// -- send the event
		ctx.Send(TodoDeleted, TodoDeletedEvent{
			Id: cmd.Id,
		})
	}
}
