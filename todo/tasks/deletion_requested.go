package tasks

import (
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/graph"
)

func onTaskDeletionRequestedLogic() []graph.Logic {
	return []graph.Logic{
		dsl.RemoveFromStore("todo", "tasks", "${! json(\"key\") }"),
		dsl.Transform(
			dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "tasks", "deleted")),
			dsl.Map("status", dsl.AsConstant(200)),
			dsl.Map("removed", dsl.ToBloblang("this")),
		),
		dsl.Catch(
			dsl.Log("ERROR", "task could not be deleted: ${!error()}"),
			dsl.Transform(
				dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "tasks", "operation_failed")),
				dsl.Map("status", dsl.AsConstant(409)),
				dsl.Map("message", dsl.AsConstant("task could not be deleted: ${!error()}")),
			),
		),
	}
}
