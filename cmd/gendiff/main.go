package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := cli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		UsageText: "gendiff [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"f"},
				DefaultText: "stylish",
				Usage:       "output formatoutput format",
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		return
	}
}
