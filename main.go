/*
Konfigurator provides a CLI tool and a library to generate Kubernetes .config files for OIDC Configured clusters.
*/
package main

import (
	"errors"
	"os"

	"github.com/MYOB-Technology/konfigurator/konfigurator"
	"github.com/urfave/cli"
)

var version = "SNAPSHOT"

func main() {
	var clientID, host, port, redirectEndpoint, ca, api, output string
	app := cli.NewApp()
	app.Name = "konfigurator"
	app.Usage = "generate a kubeconfig file with OIDC token"
	app.Version = version
	app.EnableBashCompletion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "client-id, c",
			Usage:       "oidc provider's client id (REQUIRED)",
			Destination: &clientID,
		},
		cli.StringFlag{
			Name:        "host, u",
			Usage:       "oidc provider's host url (REQUIRED)",
			Destination: &host,
		},
		cli.StringFlag{
			Name:        "api, a",
			Usage:       "url for kubernetes api (REQUIRED)",
			Destination: &api,
		},
		cli.StringFlag{
			Name:        "certificate-authority, s",
			Usage:       "certificate authority cert for kubernetes api (REQUIRED)",
			Destination: &ca,
		},
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "file to write to - defaults to stdout",
			Destination: &output,
			Value:       "",
		},
		cli.StringFlag{
			Name:        "port, p",
			Usage:       "port for the local server",
			Value:       "9000",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "redirect-endpoint, r",
			Usage:       "endpoint for oidc redirect",
			Value:       "/oauth2/callback",
			Destination: &redirectEndpoint,
		},
	}

	app.Action = func(c *cli.Context) error {

		requiredFlags := [4]string{clientID, host, api, ca}
		for _, flag := range requiredFlags {
			if flag == "" {
				return usage()
			}
		}

		konfig, err := konfigurator.NewKonfigurator(host, clientID, port, redirectEndpoint, ca, api, output)
		if err != nil {
			return err
		}

		return konfig.Orchestrate()
	}
	app.Run(os.Args)
}

func usage() error {
	return errors.New(`Missing required flag(s).

REQUIRED FLAGS:
   --client-id, -c
   --host, -u
   --api, -a
   --certificate-authority, -s
	`)
}
