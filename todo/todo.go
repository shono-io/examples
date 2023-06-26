package todo

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/local"
	"github.com/shono-io/shono/system/file"
)

func Attach(env *local.InventoryBuilder) {
	env.Scope(inventory.NewScope("todo").
		Summary("Manage todos at scale").
		Docs(`A simple and easy example application to demonstrate the power of Shono.`).
		Status(commons.StatusExperimental).
		Build())

	env.Injector(inventory.NewInjector("todo", "tasksFromFile").
		Input(file.NewInput("tasks_file", file.WithInputPath("tasks.json"))).
		OutputEvent("todo", "task", "imported").
		Logic(inventory.NewLogic().Steps(
			dsl.AsSuccessEvent(inventory.NewEventReference("todo", "task", "imported"), 200, "this")),
		).
		Build())

	AttachTask(env)
}
