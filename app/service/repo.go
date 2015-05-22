package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func pathExists(p string) bool {
	_, err := os.Stat(p)

	return err == nil
}

func cloneRepo() {
	cmd := exec.Command("git", "clone", "--branch", repoBranch, repoName, repoDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem cloning repo: %s", err)
	}
}

func pullRepo() {
	remote := fmt.Sprintf("origin/%s", repoBranch)
	cmd := exec.Command("git", "pull", repoDir, remote)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("problem pulling repo: %s", err)
	}
}

// updateRepo clones or updates the repo and returns true
// if an update occurred.
func updateRepo() {
	gitDir := filepath.Join(repoDir, ".git")

	if !pathExists(gitDir) {
		cloneRepo()
	} else {
		pullRepo()
	}

	rebuildCache(repoDir)
}

// pollRepo periodically checks the repo for updates.
func pollRepo() {
	// Check for updates every hour.
	t := time.NewTicker(interval)

	for {
		select {
		case <-t.C:
			updateRepo()
		}
	}
}