package application

import (
	"github.com/go-rabbit/amqp"
	"github.com/go-rabbit/util"
)

type (
	Application struct {
		queueName string
		rabbit    *amqp.RabbitConfig
	}

	Logger struct {
		// Stdout is true if the output needs to goto standard out
		Stdout bool `yaml:"stdout"`
		// Level is the desired log level
		Level string `yaml:"level"`
		// OutputFile is the path to the log output file
		OutputFile string `yaml:"outputFile"`
	}
)

func SetupApp() *Application {
	amqpConf := amqp.RabbitConfig{
		"localhost:32773",
		"guest",
		"guest",
	}

	logger := Logger{Stdout: true, Level: "DEBUG"}
	util.Log = logger.NewLogger()

	return &Application{
		"testing",
		&amqpConf,
	}
}
