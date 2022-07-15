package main

import (
	"github.com/gunnaringe/myuplink2mqtt/internal/server"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	var clientId string
	var clientSecret string

	var mqttServer string

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "client-id",
				Usage:       "OAuth 2.0 Client ID",
				EnvVars:     []string{"CLIENT_ID"},
				Destination: &clientId,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "client-secret",
				Usage:       "OAuth 2.0 Client secret",
				EnvVars:     []string{"CLIENT_SECRET"},
				Destination: &clientSecret,
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "homeassistant-discovery-topic",
				Usage:       "Home Assistant MQTT Discovery topic",
				EnvVars:     []string{"HOMEASSISTANT_DISCOVERY_TOPIC"},
				Required:    false,
				DefaultText: "homeassistant",
			},
			&cli.StringFlag{
				Name:        "mqtt-server",
				Usage:       "MQTT Broker URI",
				EnvVars:     []string{"MQTT_SERVER"},
				Destination: &mqttServer,
				Required:    true,
			},
		},
		Action: func(c *cli.Context) error {
			s := server.New(clientId, clientSecret, mqttServer, logger)
			return s.Run()
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal("Program exited with error", zap.Error(err))
	}
}
