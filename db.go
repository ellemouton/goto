package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type DB struct {
	path string

	contents *DBContents
	mu       sync.Mutex
}

func OpenDB() (*DB, error) {
	dir := filepath.Join(os.Getenv("HOME"), DefaultDataDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, DefaultDBName)

	if !exists(path) {
		return &DB{
			path: path,
			contents: &DBContents{
				Aliases: make(map[string]*Repo),
			},
		}, nil
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var contents DBContents
	err = json.Unmarshal(byteValue, &contents)
	if err != nil {
		return nil, err
	}

	return &DB{
		path:     path,
		contents: &contents,
	}, nil
}

func (db *DB) AddAlias(org, repo, alias string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	org = strings.ToLower(org)
	repo = strings.ToLower(repo)
	alias = strings.ToLower(alias)

	if rep, ok := db.contents.Aliases[alias]; ok {
		return fmt.Errorf("the `%s` alias has already been "+
			"registered for %s/%s", alias, rep.Org, rep.Repo)
	}

	db.contents.Aliases[alias] = &Repo{
		Org:  org,
		Repo: repo,
	}

	file, _ := json.MarshalIndent(db.contents, "", " ")

	return os.WriteFile(db.path, file, 0644)
}

func (db *DB) GetRepo(alias string) (*Repo, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	alias = strings.ToLower(alias)

	repo, ok := db.contents.Aliases[alias]
	if !ok {
		return nil, fmt.Errorf("no repo registered for this alias")
	}

	return repo, nil
}

type DBContents struct {
	Aliases map[string]*Repo `json:"aliases"`
}

type Repo struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
