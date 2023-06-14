package tasks

import (
	"github.com/shono-io/shono/dsl"
	"github.com/shono-io/shono/graph"
)

func onTaskCreationRequestedLogic() []graph.Logic {
	return []graph.Logic{
		dsl.AddToStore("todo", "task", "${! json(\"key\") }",
			dsl.Map("key", dsl.ToBloblang("this.key")),
			dsl.Map("summary", dsl.ToBloblang("this.summary")),
			dsl.Map("completed", dsl.AsConstant("false")),
			dsl.Map("timeline.createdAt", dsl.ToBloblang("@kafka_timestamp_unix")),
		),
		dsl.Transform(
			dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "task", "created")),
			dsl.Map("status", dsl.AsConstant(201)),
			dsl.Map("task", dsl.ToBloblang("this")),
		),
		dsl.Catch(
			dsl.Log("ERROR", "task could not be created: ${!error()}"),
			dsl.Transform(
				dsl.MapMeta("io_shono_kind", dsl.AsEventReference("todo", "task", "operation_failed")),
				dsl.Map("status", dsl.AsConstant(409)),
				dsl.Map("message", dsl.AsConstant("task could not be created: ${!error()}")),
			),
		),
	}
}
