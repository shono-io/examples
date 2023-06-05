package todo

import (
	"fmt"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
)

var Key = commons.NewKey("scope", "todo")

func Register(env graph.Environment) (err error) {
	err = env.RegisterScope(graph.NewScope(Key,
		graph.WithScopeName("Todos"),
		graph.WithScopeDescription("Efficient task management for the masses")))
	if err != nil {
		return fmt.Errorf("failed to register core scope: %w", err)
	}

	return nil
}
