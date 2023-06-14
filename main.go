package main

import (
	"github.com/compose-spec/compose-go/dotenv"
	"github.com/shono-io/examples/todo/tasks"
	"github.com/shono-io/shono/local"
	"github.com/shono-io/shono/runtime"
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

	// -- get the configuration file path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "shono.yaml"
	}

	// -- get the reaktors this application is responsible for
	reaktors, err := tasks.Reaktors()
	if err != nil {
		logrus.Panicf("failed to get reaktors: %v", err)
	}

	registry, err := local.Load(os.DirFS("."), configPath, local.WithReaktor(reaktors...))
	if err != nil {
		logrus.Panicf("failed to load the local registry from %q: %v", configPath, err)
	}

	if err := runtime.Run(registry); err != nil {
		logrus.Panicf("failed to run: %v", err)
	}
}
