package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "goto"
	app.Usage = "Navigate to a commit on GitHub"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		goCommand,
		registerCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
