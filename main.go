package main

import (
	"github.com/arangodb/go-driver"
	"github.com/compose-spec/compose-go/dotenv"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/shono-examples/arangodb"
	"github.com/shono-io/shono-examples/todos"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sr"
	"os"
	"strings"
)

var (
	KafkaBrokersEnv = "SHONO_KAFKA_BROKERS"
	KafkaGroupIdEnv = "SHONO_KAFKA_GROUP_ID"

	SrEndpointEnv = "SHONO_SR_ENDPOINT"

	ADBEndpointEnv = "SHONO_ARANGODB_ENDPOINT"
	ADBDatabaseEnv = "SHONO_ARANGODB_DATABASE"
	ADBUsernameEnv = "SHONO_ARANGODB_USERNAME"
	ADBPasswordEnv = "SHONO_ARANGODB_PASSWORD"

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

	// -- create the runtime
	runtime := go_shono.NewRuntime(
		go_shono.WithResource(arangodb.MustNewResource("db",
			arangodb.WithEndpoint(os.Getenv(ADBEndpointEnv)),
			arangodb.WithAuthentication(driver.BasicAuthentication(os.Getenv(ADBUsernameEnv), os.Getenv(ADBPasswordEnv))),
			arangodb.WithDatabaseName(os.Getenv(ADBDatabaseEnv)),
		)),
	)

	// -- create the agent
	agent := go_shono.NewAgent(
		"my-org",
		"my-app",
		go_shono.WithKafkaOpts(
			kgo.SeedBrokers(strings.Split(os.Getenv(KafkaBrokersEnv), ",")...),
			kgo.ConsumerGroup(os.Getenv(KafkaGroupIdEnv)),
		),
		go_shono.WithSchemaRegistryOpts(
			sr.URLs(os.Getenv(SrEndpointEnv)),
		),
		go_shono.WithReaktor(todos.Reaktors(runtime)...),
	)

	// -- run the agent
	if err := agent.Run(); err != nil {
		panic(err)
	}
}
