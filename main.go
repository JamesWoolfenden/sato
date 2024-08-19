package main

import (
	_ "embed" // required for embed
	"fmt"
	"os"
	"sort"
	"time"

	"sato/src/arm"
	"sato/src/cf"
	"sato/src/see"
	"sato/src/version"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

//goland:noinspection GoLinter
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var file string

	var destination string

	var resource string

	var flip bool

	app := &cli.App{
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:  "parse",
				Usage: "translate CFN to Terraform",
				Action: func(*cli.Context) error {
					err := cf.Parse(file, destination)
					if err != nil {
						log.Info().Msgf(err.Error())
					}

					return nil
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
					//goland:noinspection GoLinter
					fmt.Println(version.Version)

					return nil
				},
			},
			{
				Name:  "bisect",
				Usage: "translate ARM to Terraform",
				Action: func(*cli.Context) error {
					err := arm.Parse(file, destination)
					if err != nil {
						log.Info().Msgf(err.Error())
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "file",
						Aliases:     []string{"f"},
						Usage:       "Arm file to parse",
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
				Name:  "see",
				Usage: "shows equivalent Terraform resource or the reverse",
				Action: func(*cli.Context) error {
					result, err := see.Lookup(resource, flip)
					if result == nil {
						//goland:noinspection GoLinter
						return fmt.Errorf("see failed with error %w", err)
					}

					fmt.Print(*result)

					fmt.Print("\n")

					return err
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "resource",
						Aliases:     []string{"r"},
						Usage:       "target resource for conversion",
						Required:    true,
						Destination: &resource,
					},
					&cli.BoolFlag{
						Name:        "flip",
						Aliases:     []string{"f"},
						Usage:       "reverse the lookup",
						Required:    false,
						Destination: &flip,
					},
				},
			},
		},
		Name:     "sato",
		Usage:    "Translate Cloudformation or ARM to Terraform",
		Compiled: time.Time{},
		Authors:  []*cli.Author{{Name: "James Woolfenden", Email: "jim.wolf@duck.com"}},
		Version:  version.Version,
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Sato failure")
	}
}
