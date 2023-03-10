package main

import (
	"fmt"
	"github.com/urfave/cli"
)

var keywords = map[string]bool{
	"pr": true,
}

var registerCommand = cli.Command{
	Name:   "register",
	Action: registerAlias,
}

func registerAlias(ctx *cli.Context) error {
	args := ctx.Args()
	rest := args.Tail()

	if len(rest) == 0 || len(rest) > 2 {
		return fmt.Errorf("wrong number of args")
	}

	org := args.First()
	repo := args.Tail()[0]
	alias := repo
	if len(rest) == 2 {
		alias = args.Tail()[1]
	}

	if _, ok := keywords[org]; ok {
		return fmt.Errorf("org name clashes with keyword")
	}

	if _, ok := keywords[alias]; ok {
		return fmt.Errorf("alias clashes with keyword")
	}

	db, err := OpenDB()
	if err != nil {
		return err
	}

	return db.AddAlias(org, repo, alias)
}
