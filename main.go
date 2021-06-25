package main

import (
	"context"
	"fmt"
	"github.com/aljorhythm/vanilla-mock/generator"
	"github.com/aljorhythm/vanilla-mock/loader"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type cliAppArgs struct {
	dir       string
	ifaceName string
}

func newCliAppArgs(c *cli.Context) (cliAppArgs, error) {
	args := cliAppArgs{}
	args.ifaceName = c.Args().Get(0)
	var err error
	args.dir, err = os.Getwd()
	if err != nil {
		return args, err
	}
	return args, nil
}

func main() {
	app := &cli.App{
		Name:  "vanilla-mock",
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{},
		},
		Action: func(c *cli.Context) error {
			args, err := newCliAppArgs(c)
			if err != nil {
				return err
			}
			ifaceName := args.ifaceName
			dir := args.dir

			fmt.Printf("Generating mock for interface: %s\n", ifaceName)

			iface, err := loader.LoadInterface(context.Background(), dir, ifaceName)
			if err != nil {
				fmt.Errorf(err.Error())
				return err
			}

			v, err := generator.GenerateVanillaMock(iface, ifaceName)
			if err != nil {
				fmt.Errorf(err.Error())
				return err
			}

			actual := v.Output()
			fmt.Printf(actual)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
