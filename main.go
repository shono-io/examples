package main

import (
	"context"
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/go-shono/shono"
	"github.com/shono-io/go-shono/shono/local"
	"github.com/shono-io/go-shono/shono/logic"
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

	cl, err := local.NewClient("my_client", local.NewScopeRepo(), local.NewResourceRepo())
	if err != nil {
		logrus.Panicf("failed to create local client: %v", err)
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

	ctx := context.Background()

	// -- create a scope
	hr := shono.NewScope("hr", "HR Dept", "The HR Department", local.NewConceptRepo(), local.NewReaktorRepo())
	if err := cl.AddScope(ctx, hr); err != nil {
		logrus.Panicf("failed to add scope: %v", err)
	}

	// -- create a concept within the scope
	employee := shono.NewConcept("hr", "employee", "Employee", "An employee", local.NewEventRepo())
	if err := hr.AddConcept(ctx, employee); err != nil {
		logrus.Panicf("failed to add concept: %v", err)
	}

	// -- create an event within the concept
	employeeCreationRequested := shono.NewEvent("hr", "employee", "creation_requested", "Employee Creation Requested", "An employee creation was requested")
	if err := employee.AddEvent(ctx, employeeCreationRequested); err != nil {
		logrus.Panicf("failed to add event: %v", err)
	}

	employeeCreated := shono.NewEvent("hr", "employee", "created", "Employee Created", "An employee was created")
	if err := employee.AddEvent(ctx, employeeCreated); err != nil {
		logrus.Panicf("failed to add event: %v", err)
	}

	employeeCreationFailed := shono.NewEvent("hr", "employee", "creation_failed", "Employee Creation Failed", "An employee creation failed")
	if err := employee.AddEvent(ctx, employeeCreationFailed); err != nil {
		logrus.Panicf("failed to add event: %v", err)
	}

	onEmployeeCreationRequested := shono.NewReaktor("hr", "onEmployeeCreationRequested", "On Employee Creation Requested", "A reaktor that reacts to employee creation requests",
		employeeCreationRequested.Id(),
		logic.NewBenthosLogic(`
mapping: |
  root = this`),
		employeeCreated.Id(), employeeCreationFailed.Id())
	if err := hr.AddReaktor(ctx, onEmployeeCreationRequested); err != nil {
		logrus.Panicf("failed to add reaktor: %v", err)
	}

	runtime, err := shono.NewRuntime(
		shono.WithBackbone(bb),
		shono.WithReaktor(onEmployeeCreationRequested))
	if err != nil {
		logrus.Panicf("failed to create runtime: %v", err)
	}
	defer runtime.Close()

	if err := runtime.Run(ctx); err != nil {
		logrus.Panicf("failed to run: %v", err)
	}
}
