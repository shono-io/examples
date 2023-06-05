package task

import g "github.com/shono-io/shono/graph"

var onTaskCreationRequestedReaktor = g.NewReaktor(
	OnTaskCreationRequestedKey,
	CreationRequestedKey,
	g.WithReaktorDescription("Reaktor that handles task creation requests."),
	g.WithOutputEvent(CreatedKey, OperationFailedKey),
	g.WithStore(TasksStore),
	g.WithLogic(onTaskCreationRequestedLogic()...),
)

func onTaskCreationRequestedLogic() []g.Logic {
	return []g.Logic{
		g.AddToStore(TasksStoreKey, "${! json(\"key\") }",
			g.Map("key", g.ToBloblang("this.key")),
			g.Map("summary", g.ToBloblang("this.summary")),
			g.Map("completed", g.ToConstant("false")),
			g.Map("timeline.createdAt", g.ToBloblang("@kafka_timestamp_unix")),
		),
		g.Transform(
			g.MapMeta("io_shono_kind", g.ToConstant(CreatedKey.CodeString())),
			g.Map("status", g.ToConstant(201)),
			g.Map("task", g.ToBloblang("this")),
		),
		g.Catch(
			g.Log("ERROR", "task could not be created: ${!error()}"),
			g.Transform(
				g.MapMeta("io_shono_kind", g.ToConstant(OperationFailedKey.CodeString())),
				g.Map("status", g.ToConstant(409)),
				g.Map("message", g.ToConstant("task could not be created: ${!error()}")),
			),
		),
	}
}
