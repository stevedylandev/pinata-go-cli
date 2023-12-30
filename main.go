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
	Name      string            `json:"name"`
	KeyValues map[string]string `json:"keyvalues"`
}

type Region struct {
	RegionID                string `json:"regionId"`
	CurrentReplicationCount int    `json:"currentReplicationCount"`
	DesiredReplicationCount int    `json:"desiredReplicationCount"`
}

type Pin struct {
	ID            string   `json:"id"`
	IPFSPinHash   string   `json:"ipfs_pin_hash"`
	Size          int      `json:"size"`
	UserID        string   `json:"user_id"`
	DatePinned    string   `json:"date_pinned"`
	DateUnpinned  *string  `json:"date_unpinned"`
	Metadata      Metadata `json:"metadata"`
	Regions       []Region `json:"regions"`
	MimeType      string   `json:"mime_type"`
	NumberOfFiles int      `json:"number_of_files"`
}

type ListResponse struct {
	Rows []Pin `json:"rows"`
}

func main() {
	app := &cli.App{
		Name:  "pinata",
		Usage: "A CLI for uploading files to Pinata! To get started make an API key at https://app.pinata.cloud/keys, then authorize the CLI with the auth command with your JWT",
		Commands: []*cli.Command{
			{
				Name:      "auth",
				Aliases:   []string{"a"},
				Usage:     "Authorize the CLI with your Pinata JWT",
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
				Name:      "upload",
				Aliases:   []string{"u"},
				Usage:     "Upload a file or folder to Pinata",
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
			{
				Name:      "list",
				Aliases:   []string{"l"},
				Usage:     "List most recent files",
				ArgsUsage: "[number of files to return, max 1000]",
				Action: func(ctx *cli.Context) error {
					queryParam := ctx.Args().First()
					if queryParam == "" {
						queryParam = "10" // Replace with your actual default value
					}
					_, err := ListFiles(queryParam)
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
