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
	Name      string                 `json:"name"`
	KeyValues map[string]interface{} `json:"keyvalues"`
}

type Pin struct {
	ID            string   `json:"id"`
	IPFSPinHash   string   `json:"ipfs_pin_hash"`
	Size          int      `json:"size"`
	UserID        string   `json:"user_id"`
	DatePinned    string   `json:"date_pinned"`
	DateUnpinned  *string  `json:"date_unpinned"`
	Metadata      Metadata `json:"metadata"`
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
				Name:      "pin",
				Aliases:   []string{"p"},
				Usage:     "Pin an existing CID on IPFS to Pinata",
				ArgsUsage: "[CID of file on IPFS]",
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
				Name:      "delete",
				Aliases:   []string{"d"},
				Usage:     "Delete a file by CID",
				ArgsUsage: "[CID of file]",
				Action: func(ctx *cli.Context) error {
					cid := ctx.Args().First()
					if cid == "" {
						return errors.New("no CID provided")
					}
					err := Delete(cid)
					return err
				},
			},
			{
				Name:      "list",
				Aliases:   []string{"l"},
				Usage:     "List most recent files",
				ArgsUsage: "[List your most recent files]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cid",
						Aliases: []string{"c"},
						Value:   "null",
						Usage:   "Search files by CID",
					},
					&cli.StringFlag{
						Name:    "amount",
						Aliases: []string{"a"},
						Value:   "10",
						Usage:   "The number of files you would like to return, default 10 max 1000",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "null",
						Usage:   "The name of the file",
					},
					&cli.StringFlag{
						Name:    "status",
						Aliases: []string{"s"},
						Value:   "pinned",
						Usage:   "Status of the file. Options are 'pinned', 'unpinned', or 'all'. Default: 'pinned'",
					},
					&cli.StringFlag{
						Name:    "pageOffset",
						Aliases: []string{"p"},
						Value:   "null",
						Usage:   "Allows you to paginate through files. If your file amount is 10, then you could set the pageOffset to '10' to see the next 10 files.",
					},
				},
				Action: func(ctx *cli.Context) error {
					cid := ctx.String("cid")
					amount := ctx.String("amount")
					name := ctx.String("name")
					status := ctx.String("status")
					offset := ctx.String("pageOffset")
					_, err := ListFiles(amount, cid, name, status, offset)
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
