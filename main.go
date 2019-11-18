package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sebcoetzee/cloudflare-dynamic-dns/app"
	"github.com/sebcoetzee/cloudflare-dynamic-dns/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	setupLogrus()
	var config config.Config

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Print(c.App.Version)
	}

	app := &cli.App{
		Version: "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "api_token",
				Required:    true,
				Usage:       "Cloudflare API Token",
				Destination: &config.APIToken,
			},
			&cli.StringFlag{
				Name:        "zone_name",
				Required:    true,
				Usage:       "Cloudflare Zone name that will be used e.g. 'example.com'",
				Destination: &config.ZoneName,
			},
			&cli.StringFlag{
				Name:        "subdomain",
				Required:    true,
				Usage:       "Cloudflare subdomain that will be used to set the DNS record e.g. 'home' which will be appended to the zone name to form 'home.example.com'",
				Destination: &config.Subdomain,
			},
		},
		Action: func(c *cli.Context) error {
			var err error
			ticker := time.NewTicker(5 * 60 * time.Second)
			done := make(chan bool)

			app := app.NewApp(config)

			go func() {
				err = app.UpdateDNS()
				for {
					select {
					case <-ticker.C:
						err = app.UpdateDNS()
					case <-done:
						os.Exit(0)
					}
				}
			}()

			sigChan := make(chan os.Signal)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			<-sigChan
			done <- true
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setupLogrus() {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
}
