package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/shono/cloud"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"net"
	"os"
	"strings"
	"time"
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

	tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
	cc, err := cloud.NewConfluentClient(
		os.Getenv(ConfluentEnvironmentIdEnv),
		os.Getenv(ConfluentClusterIdEnv),
		os.Getenv(ConfluentClusterAPIEndpointEnv),
		os.Getenv(ConfluentApiKeyEnv),
		os.Getenv(ConfluentApiSecretEnv),
		kgo.SeedBrokers(strings.Split(os.Getenv(KafkaBrokersEnv), ",")...),
		kgo.SASL(plain.Auth{User: os.Getenv(KafkaApiKeyEnv), Pass: os.Getenv(KafkaApiSecretEnv)}.AsMechanism()),
		kgo.Dialer(tlsDialer.DialContext),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create confluent client: %w", err))
	}

	err = cc.CreateTopicAcl(context.Background(), "ddd", "sa-d1ndkd")
	if err != nil {
		panic(err)
	}
}
