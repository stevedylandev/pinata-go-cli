package main

import (
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
