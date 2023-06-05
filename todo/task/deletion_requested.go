package task

import g "github.com/shono-io/shono/graph"

var onTaskDeletionRequestedReaktor = g.NewReaktor(
	OnTaskDeletionRequestedKey,
	DeletionRequestedKey,
	g.WithReaktorDescription("Reaktor that handles task deletion requests."),
	g.WithOutputEvent(DeletedKey, OperationFailedKey),
	g.WithStore(TasksStore),
	g.WithLogic(onTaskDeletionRequestedLogic()...),
)

func onTaskDeletionRequestedLogic() []g.Logic {
	return []g.Logic{
		g.RemoveFromStore(TasksStoreKey, "${! json(\"key\") }"),
		g.Transform(
			g.MapMeta("io_shono_kind", g.ToConstant(DeletedKey.CodeString())),
			g.Map("status", g.ToConstant(200)),
			g.Map("removed", g.ToBloblang("this")),
		),
		g.Catch(
			g.Log("ERROR", "task could not be deleted: ${!error()}"),
			g.Transform(
				g.MapMeta("io_shono_kind", g.ToConstant(OperationFailedKey.CodeString())),
				g.Map("status", g.ToConstant(409)),
				g.Map("message", g.ToConstant("task could not be deleted: ${!error()}")),
			),
		),
	}
}
