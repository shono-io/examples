package main

import (
	"github.com/arangodb/go-driver"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/shono-examples/arangodb"
	"github.com/shono-io/shono-examples/todos"
	"github.com/twmb/franz-go/pkg/kgo"
	"os"
	"strings"
)

var (
	KafkaBrokersEnv = "SHONO_KAFKA_BROKERS"
	KafkaGroupIdEnv = "SHONO_KAFKA_GROUP_ID"

	ADBEndpointEnv = "SHONO_ARANGODB_ENDPOINT"
	ADBDatabaseEnv = "SHONO_ARANGODB_DATABASE"
	ADBUsernameEnv = "SHONO_ARANGODB_USERNAME"
	ADBPasswordEnv = "SHONO_ARANGODB_PASSWORD"
)

func main() {
	// -- create the agent
	agent := go_shono.NewAgent(
		"",
		"",
		go_shono.WithKafkaOpts(
			kgo.SeedBrokers(strings.Split(os.Getenv(KafkaBrokersEnv), ",")...),
			kgo.ConsumerGroup(os.Getenv(KafkaGroupIdEnv)),
		),
		go_shono.WithResource(arangodb.MustNewResource("db",
			arangodb.WithEndpoint(os.Getenv(ADBEndpointEnv)),
			arangodb.WithAuthentication(driver.BasicAuthentication(os.Getenv(ADBUsernameEnv), os.Getenv(ADBPasswordEnv))),
			arangodb.WithDatabaseName(os.Getenv(ADBDatabaseEnv)),
		)),
		go_shono.WithReaktor(todos.Reaktors()...),
	)

	// -- run the agent
	if err := agent.Run(); err != nil {
		panic(err)
	}
}
