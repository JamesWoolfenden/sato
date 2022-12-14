package main

import (
	_ "embed" //required for embed
	"fmt"
	"os"
	sato "sato/src"
	"sort"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var file string
	var destination string

	app := &cli.App{
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:  "parse",
				Usage: "translate CFN to Terraform",
				Action: func(*cli.Context) error {
					err := sato.Parse(file, destination)
					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Usage:       "Cloudformation file to parse",
						Required:    true,
						Destination: &file,
					},
					&cli.StringFlag{
						Name:        "destination",
						Aliases:     []string{"d"},
						Usage:       "Destination to write Terraform",
						Value:       ".sato",
						Destination: &destination,
					},
				},
			},
			{
				Name:      "version",
				Aliases:   []string{"v"},
				Usage:     "Outputs the application version",
				UsageText: "sato version",
				Action: func(*cli.Context) error {
					fmt.Println(sato.Version)
					return nil
				},
			},
		},
		Name:     "sato",
		Usage:    "Translate Cloudformation to Terraform",
		Compiled: time.Time{},
		Authors:  []*cli.Author{{Name: "James Woolfenden", Email: "support@bridgecrew.io"}},
		Version:  sato.Version,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Sato failure")
	}

}
