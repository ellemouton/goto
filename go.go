package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"os/exec"
)

const (
	DefaultDataDir = ".goto"
	DefaultDBName  = "db.json"
)

var goCommand = cli.Command{
	Name:   "go",
	Action: goToWebPage,
	Subcommands: []cli.Command{
		goToPR,
	},
}

var goToPR = cli.Command{
	Name:   "pr",
	Action: gotoPR,
}

func gotoPR(ctx *cli.Context) error {
	args := ctx.Args()

	if !args.Present() {
		return fmt.Errorf("not enough arguments")
	}

	switch len(args.Tail()) {
	case 1:
		return handleRegisteredRepo(args, true)
	case 2:
		return handleUnregisteredRepo(args)
	default:
		return fmt.Errorf("invalid number of args")
	}
}

func goToWebPage(ctx *cli.Context) error {
	args := ctx.Args()

	if !args.Present() {
		return fmt.Errorf("not enough arguments")
	}

	switch len(args.Tail()) {
	case 1:
		return handleRegisteredRepo(args, false)
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

func handleRegisteredRepo(args cli.Args, pr bool) error {
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

	commitHash := args.Tail()[0]

	if pr {
		return openPRURL(repo.Org, repo.Repo, commitHash)
	}

	return openCommitURL(repo.Org, repo.Repo, commitHash)
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

type GHPRs []*GHPRInfo

type GHPRInfo struct {
	URL string `json:"html_url"`
}

func openPRURL(org, repo, commit string) error {
	queryURI := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/commits/%s/pulls", org,
		repo, commit,
	)

	resp, err := http.Get(queryURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var prs GHPRs
	err = json.NewDecoder(resp.Body).Decode(&prs)
	if err != nil {
		return err
	}

	switch len(prs) {
	case 0:
		return fmt.Errorf("no PRs associated with this commit")
	case 1:
	default:
		var links string
		for _, link := range prs {
			links += fmt.Sprintf("%s\n", link.URL)
		}

		fmt.Println("There are multiple PRs associated with this "+
			"commit. Here are the links: ", links)

		return nil
	}

	return exec.Command("open", prs[0].URL).Run()
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
