package main

import (
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/shono"
	"github.com/shono-io/shono/backbone"
	"github.com/shono-io/shono/benthos"
	"github.com/shono-io/shono/graph"
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

	bb := backbone.NewKafkaBackbone(map[string]any{
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
	}, backbone.PerScopeLogStrategy)

	s := core.CoreScope

	gen := benthos.NewGenerator("cloud", bb, 2)
	res, err := gen.Generate(s)
	if err != nil {
		logrus.Panicf("failed to generate: %v", err)
	}

	logrus.Infof("configuration:")
	if err := res.Write(os.Stdout); err != nil {
		logrus.Panicf("failed to write configuration: %v", err)
	}

	// -- create the output directory if it doesn't exist
	if err := os.MkdirAll("out", os.ModePerm); err != nil {
		logrus.Panicf("failed to create output directory: %v", err)
	}

	// -- write the configuration to a file
	f, err := os.Create("out/benthos.yaml")
	if err != nil {
		logrus.Panicf("failed to create output file: %v", err)
	}
	defer f.Close()

	if err := res.Write(f, shono.NonSecure()); err != nil {
		logrus.Panicf("failed to write configuration: %v", err)
	}
}

func buildEnvironment() graph.Environment {

}
