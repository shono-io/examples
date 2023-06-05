package task

import g "github.com/shono-io/shono/graph"

var onTaskCompletionRequestedReaktor = g.NewReaktor(
	OnTaskCompletionRequestedKey,
	CompletionRequestedKey,
	g.WithReaktorDescription("Reaktor that handles task completion requests."),
	g.WithOutputEvent(FinishedKey, OperationFailedKey),
	g.WithStore(TasksStore),
	g.WithLogic(onTaskCompletionRequestedLogic()...),
)

func onTaskCompletionRequestedLogic() []g.Logic {
	return []g.Logic{
		g.GetFromStore(TasksStoreKey, "${! json(\"key\") }"),
		g.SetInStore(TasksStoreKey, "${! json(\"key\") }",
			g.MapRoot(),
			g.Map("finished", g.ToConstant(true)),
		),
		g.Transform(
			g.MapMeta("io_shono_kind", g.ToConstant(FinishedKey.CodeString())),
			g.MapRoot(),
		),
		g.Catch(
			g.Log("ERROR", "task could not be completed: ${!error()}"),
			g.Transform(
				g.MapMeta("io_shono_kind", g.ToConstant(OperationFailedKey.CodeString())),
				g.Map("status", g.ToConstant(409)),
				g.Map("message", g.ToConstant("task could not be created: ${!error()}")),
			),
		),
	}
}
