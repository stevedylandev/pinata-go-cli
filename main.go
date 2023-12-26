package main

import (
	"errors"
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
type Metadata struct {
	Name string `json:"name"`
}

func main() {
	app := &cli.App{
		Name:  "pinata",
		Usage: "A CLI for uploading files to Pinata! To get started make an API key at https://app.pinata.cloud/keys, then authorize the CLI with the auth command with your JWT",
		Commands: []*cli.Command{
			{
				Name:    "auth",
				Aliases: []string{"a"},
				Usage:   "Authorize the CLI with your Pinata JWT",
        ArgsUsage: "[your Pinata JWT]",
				Action: func(ctx *cli.Context) error {
					jwt := ctx.Args().First()
					if jwt == "" {
						return errors.New("no jwt supplied")
					}
					err := SaveJWT(jwt)
					return err
				},
			},
			{
				Name:    "upload",
				Aliases: []string{"u"},
				Usage:   "Upload a file or folder to Pinata",
        ArgsUsage: "[path to file]",
				Action: func(ctx *cli.Context) error {
					filePath := ctx.Args().First()
					if filePath == "" {
						return errors.New("no file path supplied")
					}
					_, err := Upload(filePath)
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
