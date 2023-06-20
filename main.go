package main

import (
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/examples/todo"
	"github.com/shono-io/shono/artifacts/benthos"
	"github.com/shono-io/shono/inventory"
	"github.com/shono-io/shono/local"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	LogLevelEnv = "LOG_LEVEL"
)

func main() {
	if err := dotenv.Load(); err != nil {
		logrus.Panicf("failed to load .env file: %v", err)
	}

	ll := os.Getenv(LogLevelEnv)
	if ll != "" {
		lv, err := logrus.ParseLevel(ll)
		if err != nil {
			logrus.Panicf("failed to parse log level: %v", err)
		} else {
			logrus.SetLevel(lv)
		}
	}

	ib := local.NewInventory()
	todo.Attach(ib)
	inv := ib.Build()

	// -- generate the artifacts for all the reaktors in the registry
	artifact, err := benthos.NewConceptGenerator().Generate(inv, inventory.NewConceptReference("todo", "task"))
	if err != nil {
		logrus.Panicf("failed to generate concept artifact: %v", err)
	}

	if err := local.DumpArtifact("artifacts", artifact); err != nil {
		logrus.Panicf("failed to dump artifact: %v", err)
	}
}
