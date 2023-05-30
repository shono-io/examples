package main

import (
	"context"
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/shono"
	"github.com/shono-io/shono/logic"
	"github.com/sirupsen/logrus"
	"os"
)

var (
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

	bb := shono.NewKafkaBackbone(map[string]any{
		"seed_brokers": []string{
			os.Getenv(KafkaBrokersEnv),
		},
		"tls": map[string]any{
			"enabled": true,
		},
		"sasl": []map[string]any{
			{
				"mechanism": "PLAIN",
				"username":  os.Getenv(KafkaApiKeyEnv),
				"password":  os.Getenv(KafkaApiSecretEnv),
			},
		},
	}, shono.PerScopeLogStrategy)

	// -- create events
	var (
		employeeCreationRequested = shono.NewEvent("hr", "employee", "creation_requested")
		employeeCreated           = shono.NewEvent("hr", "employee", "created")
		employeeCreationFailed    = shono.NewEvent("hr", "employee", "creation_failed")

		employeeDeletionRequested = shono.NewEvent("hr", "employee", "deletion_requested")
		employeeDeleted           = shono.NewEvent("hr", "employee", "deleted")
		employeeDeletionFailed    = shono.NewEvent("hr", "employee", "deletion_failed")
	)

	// -- create a first reaktor that listens to the employee creation requested event
	onEmployeeCreationRequested := shono.NewReaktor("hr", "onEmployeeCreationRequested",
		employeeCreationRequested.Id(),
		logic.NewBenthosLogic(`mapping: root = this`),
		shono.WithOutputEvent(employeeCreated.Id()),
		shono.WithOutputEvent(employeeCreationFailed.Id()))

	// -- create a second reaktor that listens to the employee deletion requested event
	onEmployeeDeletionRequested := shono.NewReaktor("hr", "onEmployeeDeletionRequested",
		employeeDeletionRequested.Id(),
		logic.NewBenthosLogic(`mapping: root = this`),
		shono.WithOutputEvent(employeeDeleted.Id()),
		shono.WithOutputEvent(employeeDeletionFailed.Id()))

	// -- create a runtime for both reaktors
	runtime, err := shono.NewRuntime(
		shono.WithBackbone(bb),
		shono.WithReaktor(onEmployeeCreationRequested),
		shono.WithReaktor(onEmployeeDeletionRequested))
	if err != nil {
		logrus.Panicf("failed to create runtime: %v", err)
	}
	defer runtime.Close()

	// -- execute the runtime
	if err := runtime.Run(context.Background()); err != nil {
		logrus.Panicf("failed to run: %v", err)
	}
}
