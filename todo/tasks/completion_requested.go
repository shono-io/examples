package tasks

import (
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/graph"
)

func onTaskCompletionRequestedLogic() []graph.Logic {
	return []graph.Logic{
		dsl.GetFromStore("todo", "task", "${! json(\"key\") }"),
		dsl.SetInStore("todo", "task", "${! json(\"key\") }",
			dsl.MapRoot(),
			dsl.Map("finished", dsl.AsConstant(true)),
		),
		dsl.Transform(
			dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "task", "finished")),
			dsl.MapRoot(),
		),
		dsl.Catch(
			dsl.Log("ERROR", "task could not be completed: ${!error()}"),
			dsl.Transform(
				dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "task", "operation_failed")),
				dsl.Map("status", dsl.AsConstant(409)),
				dsl.Map("message", dsl.AsConstant("task could not be created: ${!error()}")),
			),
		),
	}
}
