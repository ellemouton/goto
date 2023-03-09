package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os/exec"
)

const (
	DefaultDataDir = ".goto"
	DefaultDBName  = "db.json"
)

var goCommand = cli.Command{
	Name:   "go",
	Action: goToWebPage,
}

func goToWebPage(ctx *cli.Context) error {
	args := ctx.Args()

	if !args.Present() {
		return fmt.Errorf("not enough arguments")
	}

	switch len(args.Tail()) {
	case 1:
		return handleRegisteredRepo(args)
	case 2:
		return handleUnregisteredRepo(args)
	default:
		return goToRepo(args.First())
	}
}

func goToRepo(alias string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}

	repo, err := db.GetRepo(alias)
	if err != nil {
		return err
	}

	return openRepoURL(repo.Org, repo.Repo)
}

func handleRegisteredRepo(args cli.Args) error {
	if len(args.Tail()) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	db, err := OpenDB()
	if err != nil {
		return err
	}

	repo, err := db.GetRepo(args.First())
	if err != nil {
		return err
	}

	return openCommitURL(repo.Org, repo.Repo, args.Tail()[0])
}

func handleUnregisteredRepo(args cli.Args) error {
	if len(args.Tail()) != 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	org := args.First()
	rest := args.Tail()

	repo := rest[0]
	commit := rest[1]

	return openCommitURL(org, repo, commit)
}

func openCommitURL(org, repo, commit string) error {
	url := fmt.Sprintf(
		"https://github.com/%s/%s/commit/%s", org, repo, commit,
	)

	return exec.Command("open", url).Run()
}

func openRepoURL(org, repo string) error {
	url := fmt.Sprintf("https://github.com/%s/%s", org, repo)

	return exec.Command("open", url).Run()
}
