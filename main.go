package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type ResponseData struct {
	IpfsHash    string `json:"IpfsHash"`
	PinSize     int    `json:"PinSize"`
	Timestamp   string `json:"Timestamp"`
	IsDuplicate bool   `json:"isDuplicate"`
}

type Options struct {
	CidVersion int `json:"cidVersion"`
}
type KeyValues struct {
	WhimseyLevel int `json:"whimsey_level"`
}
type Metadata struct {
	Name      string    `json:"name"`
	Keyvalues KeyValues `json:"keyvalues"`
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "auth",
				Aliases: []string{"a"},
				Usage:   "Authorize the CLI with your Pinata JWT",
				Action: func(cCtx *cli.Context) error {
					SaveJWT(string(cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Upload a file or folder to Pinata",
				Action: func(cCtx *cli.Context) error {
					Upload(string(cCtx.Args().First()))
					return nil
				},
			},
			{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "options for task templates",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a new template",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("new task template: ", cCtx.Args().First())
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove an existing template",
						Action: func(cCtx *cli.Context) error {
							fmt.Println("removed task template: ", cCtx.Args().First())
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
