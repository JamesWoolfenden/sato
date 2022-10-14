package main

import (
	_ "embed" //required for embed
	"github.com/urfave/cli/v2"
	"log"
	"os"
	sato "sato/src"
)

func main() {
	var file string
	var destination string

	app := &cli.App{
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Name:                 "sato",
		Usage:                "translate CFN to Terraform",
		Commands: []*cli.Command{
			{
				Name:  "parse",
				Usage: "translate CFN to Terraform",
				Action: func(*cli.Context) error {
					sato.Parse(file, destination)
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
