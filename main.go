package main

import (
	"fmt"
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/examples/todo"
	"github.com/shono-io/examples/todo/task"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
	"github.com/shono-io/shono/local"
	"github.com/shono-io/shono/runtime"
	"github.com/shono-io/shono/systems/backbone"
	"github.com/shono-io/shono/systems/storage/arangodb"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	ADBEndpointEnv = "ARANGODB_ENDPOINT"
	ADBDatabaseEnv = "ARANGODB_DATABASE"
	ADBUsernameEnv = "ARANGODB_USERNAME"
	ADBPasswordEnv = "ARANGODB_PASSWORD"

	KafkaBrokersEnv   = "KAFKA_BROKERS"
	KafkaApiKeyEnv    = "KAFKA_API_KEY"
	KafkaApiSecretEnv = "KAFKA_API_SECRET"

	ConfluentEnvironmentIdEnv      = "CONFLUENT_ENVIRONMENT_ID"
	ConfluentClusterIdEnv          = "CONFLUENT_CLUSTER_ID"
	ConfluentClusterAPIEndpointEnv = "CONFLUENT_CLUSTER_API_ENDPOINT"
	ConfluentApiKeyEnv             = "CONFLUENT_API_KEY"
	ConfluentApiSecretEnv          = "CONFLUENT_API_SECRET"

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

	bb := backbone.NewKafkaBackbone(backbone.KafkaBackboneConfig{
		BootstrapServers: []string{os.Getenv(KafkaBrokersEnv)},
		TLS: &backbone.TLS{
			Enabled: true,
		},
		SASL: []backbone.SASLConfig{
			{
				Mechanism: "PLAIN",
				Username:  os.Getenv(KafkaApiKeyEnv),
				Password:  os.Getenv(KafkaApiSecretEnv),
			},
		},
		LogStrategy: backbone.PerScopeLogStrategy,
	})

	env := local.NewEnvironment(bb)
	if err := register(env); err != nil {
		logrus.Panicf("failed to register: %v", err)
	}

	if err := runtime.Run(env); err != nil {
		logrus.Panicf("failed to run: %v", err)
	}
}

func register(env graph.Environment) error {
	// -- register the arangodb storage
	adbStorage, err := arangodb.NewStorage(
		commons.NewKey("storage", "arangodb"),
		[]string{os.Getenv(ADBEndpointEnv)},
		os.Getenv(ADBUsernameEnv),
		os.Getenv(ADBPasswordEnv),
		os.Getenv(ADBDatabaseEnv),
	)
	if err != nil {
		return fmt.Errorf("failed to create arangodb storage: %w", err)
	}
	if err := env.RegisterStorage(*adbStorage); err != nil {
		return fmt.Errorf("failed to register arangodb storage: %w", err)
	}

	if err := todo.Register(env); err != nil {
		return err
	}

	if err := task.Register(env); err != nil {
		return fmt.Errorf("failed to register task: %w", err)
	}

	return nil
}
