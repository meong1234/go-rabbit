package main

import (
	"github.com/go-rabbit/application"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := application.SetupApp()

	clientApp := cli.NewApp()
	clientApp.Name = "go-rabbit"
	clientApp.Version = "0.0.1"
	clientApp.Commands = []cli.Command{
		{
			Name:        "publish",
			Description: "start publisher",
			Action: func(c *cli.Context) error {
				daemon := app.NewPublisherDaemon()
				return application.AppRunner(daemon)
			},
		},
		{
			Name:        "subscribe",
			Description: "start subscriber",
			Action: func(c *cli.Context) error {
				daemon := app.NewSubscriberDaemon()
				return application.AppRunner(daemon)
			},
		},
	}

	clientApp.Run(os.Args)
}
