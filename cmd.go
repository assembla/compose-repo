package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

func statusCmd(c *cli.Context) {
	repos := unmarshalRepos()

	for _, rep := range repos {
		log.Printf("%s", rep.Path)
		rep.Status()
	}
}

func initCmd(c *cli.Context) {
	repoURL := c.String("repo")
	if repoURL == "" {
		log.Fatal("please provide --repo, -r param")
	}

	clone(repoURL, ComposerRepo)

	if err := os.Symlink(filepath.Join(ComposerRepo, ComposeFile), ComposeFile); err != nil {
		log.Fatal(err)
	}
}

func syncCmd(c *cli.Context) {
	fetch(ComposerRepo)
	repos := unmarshalRepos()

	for i, rep := range repos {
		log.Printf("syncing %s (%d/%d)\n", rep.Path, i+1, len(repos))
		rep.Get()
	}
}
