package main

import (
	"context"
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/shono"
	"github.com/shono-io/shono/benthos"
	"github.com/shono-io/shono/logic"
	"github.com/shono-io/shono/store"
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

	// -- create the backbone we will use for the reaktor
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

	hrScope := shono.NewScope("hr")
	employeeConcept := hrScope.NewConcept("employee")

	// -- create the store we will use for employee concepts
	employeeStore := store.NewArangodbStore(employeeConcept, "employee_store",
		os.Getenv(ADBEndpointEnv), os.Getenv(ADBDatabaseEnv), "employees",
		os.Getenv(ADBUsernameEnv), os.Getenv(ADBPasswordEnv))

	// -- create events
	var (
		employeeCreationRequested = shono.NewEvent(employeeConcept.Key(), "creation_requested")
		employeeCreated           = shono.NewEvent(employeeConcept.Key(), "created")
		employeeCreationFailed    = shono.NewEvent(employeeConcept.Key(), "creation_failed")

		employeeDeletionRequested = shono.NewEvent(employeeConcept.Key(), "deletion_requested")
		employeeDeleted           = shono.NewEvent(employeeConcept.Key(), "deleted")
		employeeDeletionFailed    = shono.NewEvent(employeeConcept.Key(), "deletion_failed")
	)

	// -- create a first reaktor that listens to the employee creation requested event
	onEmployeeCreationRequested := shono.NewReaktor(hrScope.Key(), "onEmployeeCreationRequested",
		employeeCreationRequested.Id(),
		logic.NewBenthosLogic(`mapping: root = this`),
		shono.WithOutputEvent(employeeCreated.Id()),
		shono.WithOutputEvent(employeeCreationFailed.Id()),
		shono.WithStore(employeeStore))

	// -- create a second reaktor that listens to the employee deletion requested event
	onEmployeeDeletionRequested := shono.NewReaktor(hrScope.Key(), "onEmployeeDeletionRequested",
		employeeDeletionRequested.Id(),
		logic.NewBenthosLogic(`mapping: root = this`),
		shono.WithOutputEvent(employeeDeleted.Id()),
		shono.WithOutputEvent(employeeDeletionFailed.Id()),
		shono.WithStore(employeeStore))

	// -- create a runtime for both reaktors
	runtime, err := benthos.NewRuntime(
		benthos.WithBackbone(bb),
		benthos.WithReaktor(onEmployeeCreationRequested),
		benthos.WithReaktor(onEmployeeDeletionRequested))
	if err != nil {
		logrus.Panicf("failed to create runtime: %v", err)
	}
	defer runtime.Close()

	// -- execute the runtime
	if err := runtime.Run(context.Background()); err != nil {
		logrus.Panicf("failed to run: %v", err)
	}
}
