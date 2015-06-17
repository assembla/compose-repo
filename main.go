package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
)

const (
	ComposerRepo  = ".compose-repo"
	Manifest      = "manifest.json"
	ComposeFile   = "docker-compose.yml"
	defaultBranch = "master"
	defaultRemote = "origin/master"
)

type Repo struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

func (r Repo) Get() error {
	if r.exists() {
		fetch(r.Path)
	} else {
		clone(r.URL, r.Path)
	}
	return nil
}

func (r Repo) exists() bool {
	f, err := os.Open(r.Path)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Stat()
	if err != nil {
		return false
	}

	return true
}

var repos []Repo

func Git(params ...string) *exec.Cmd {
	cmd := exec.Command("git", params...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd
}

func GitInDir(dir string, params ...string) *exec.Cmd {
	cmd := Git(params...)
	cmd.Dir = dir
	return cmd
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("===> ")
	app := cli.NewApp()
	app.Name = "composer-repo"
	app.Usage = "sync multiple git repos"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:        "init",
			Aliases:     []string{"i"},
			Usage:       "initialize repo with manifest.json and docker-compose.yml",
			Description: "",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "repo,r",
					Usage: "path to manifest repo",
				},
			},
			Action: initCmd,
		},
		{
			Name:        "sync",
			Aliases:     []string{"s"},
			Usage:       "syncronize repos from manifest.json",
			Description: "",
			Action:      syncCmd,
		},
		{
			Name:        "status",
			Aliases:     []string{"st"},
			Usage:       "show status across all repos",
			Description: "",
			Action:      statusCmd,
		},
	}

	app.Run(os.Args)
}

func statusCmd(c *cli.Context) {
	// TODO: show git status with appropriate format from all repos at once
	// TODO: check if dir is not in manifest
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

	f, err := os.Open(filepath.Join(ComposerRepo, Manifest))

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)

	if err := json.Unmarshal(data, &repos); err != nil {
		log.Fatal(err)
	}

	for i, rep := range repos {
		log.Printf("syncing %s (%d/%d)\n", rep.Path, i+1, len(repos))
		rep.Get()
	}
}

func clone(src, dst string) {
	if err := Git("clone", src, dst).Run(); err != nil {
		log.Fatal(err)
	}
}

func fetch(dir string) {
	if err := GitInDir(dir, "fetch", "origin").Run(); err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer

	cmd := GitInDir(dir, "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = &b
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	br := b.Bytes()
	curBranch := string(bytes.TrimSpace(br))

	if curBranch != defaultBranch {
		log.Printf("found that current branch is %q, switching to %q", curBranch, defaultBranch)
		if err := GitInDir(dir, "checkout", defaultBranch).Run(); err != nil {
			log.Fatal(err)
		}
	}

	args := []string{"reset", "--hard", defaultRemote}
	if err := GitInDir(dir, args...).Run(); err != nil {
		log.Fatal(err)
	}

	if curBranch != defaultBranch {
		log.Printf("switch back to %q branch after reseting %q to %q", curBranch, defaultBranch, defaultRemote)
		if err := GitInDir(dir, "checkout", curBranch).Run(); err != nil {
			log.Fatal(err)
		}
	}
}
