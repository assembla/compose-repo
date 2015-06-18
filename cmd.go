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
		rep.Status().Run()
	}
}

func initCmd(c *cli.Context) {
	repoURL := c.String("repo")
	if repoURL == "" {
		log.Fatal("please provide --repo, -r param")
	}

	repo := Repo{Path: ComposerRepo, URL: repoURL}
	repo.Clone().Run()

	if err := os.Symlink(filepath.Join(ComposerRepo, ComposeFile), ComposeFile); err != nil {
		log.Fatal(err)
	}
}

func syncCmd(c *cli.Context) {
	repo := Repo{Path: ComposerRepo}
	repo.RemoteUpdate("-p").Reset(defaultRemote).Run()

	repos := unmarshalRepos()
	for i, rep := range repos {
		log.Printf("syncing %s (%d/%d)\n", rep.Path, i+1, len(repos))

		if RepoExists(rep) {
			remoteArgs := []string{}
			if c.Bool("prune") {
				remoteArgs = append(remoteArgs, "-p")
			}

			rep = rep.RemoteUpdate(remoteArgs...)

			if c.Bool("hard") {
				rep = rep.Reset(defaultRemote)
			}
		} else {
			rep = rep.Clone()
		}
		rep.Run()
	}
}
