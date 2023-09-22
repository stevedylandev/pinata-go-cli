package main

import (
	"errors"
	"log"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

var primaryStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.CompleteColor{TrueColor: "#8A79FF", ANSI256: "99", ANSI: "99"})

var successStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.CompleteColor{TrueColor: "#00CC92", ANSI256: "42", ANSI: "42"})

var failureStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.CompleteColor{TrueColor: "#F04438", ANSI256: "166", ANSI: "166"})

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
