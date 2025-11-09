package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"code"
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
		Action: func(ctx context.Context, c *cli.Command) error {
			fmt.Println()
			if c.NArg() != 2 {
				return fmt.Errorf("need exactly 2 file paths")
			}

			filepath1 := c.Args().Get(0)
			filepath2 := c.Args().Get(1)

			result, err := code.GenDiff(filepath1, filepath2)
			if err != nil {
				return err
			}
			fmt.Println(result)
			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
