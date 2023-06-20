package todo

import (
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/local"
)

func Attach(env *local.EnvironmentBuilder) {
	env.Scope(inventory.NewScope("todo").
		Summary("Manage todos at scale").
		Docs(`A simple and easy example application to demonstrate the power of Shono.`).
		Status(commons.StatusExperimental).
		Build())

	//env.Injector(inventory.NewInjector("tasksFromFile").
	//	Input(file.NewInput(file.WithInputPath("tasks.json"))).
	//	OutputEvent("todo", "task", "imported").
	//	Logic(inventory.NewLogic().Steps(
	//		dsl.AsEvent(inventory.NewEventReference("todo", "task", "imported"))),
	//	).
	//	Build())

	AttachTask(env)
}
