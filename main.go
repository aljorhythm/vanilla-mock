package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aljorhythm/vanilla-mock/common"
	"github.com/aljorhythm/vanilla-mock/generator"
	"github.com/aljorhythm/vanilla-mock/loader"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

func loadAndGenerate(dir string, ifaceName string) (string, error) {
	if ifaceName == "" {
		return "", errors.New("empty interface name")
	}
	iface, err := loader.LoadInterface(context.Background(), dir, ifaceName)
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}

	v, err := generator.GenerateVanillaMock(iface, ifaceName)
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}

	actual := v.Output()
	fmt.Printf(actual)

	return actual, err
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

			mockStructString, err := loadAndGenerate(dir, ifaceName)

			if err != nil {
				fmt.Printf("warning: error generating %s", err.Error())
				return err
			}

			outputDir := filepath.Join(dir, "mock")

			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				// path/to/whatever does not exist

				err := os.Mkdir(outputDir, os.ModeDir)
				if err != nil {
					fmt.Printf("warning: cannot create directory %s\n", outputDir)
				}
			}

			newFilename := fmt.Sprintf("%s.go", common.ToSnakeCase(ifaceName))
			newFilepath := filepath.Join(outputDir, newFilename)
			fmt.Printf("writing to file %s\n", newFilepath)

			err = ioutil.WriteFile(newFilepath, []byte(mockStructString), 0644)

			if err != nil {
				fmt.Printf("warning: failed writing to file %s\n", newFilepath)
			}

			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
