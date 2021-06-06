package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
)

func cloneAction(actionID string) (string, func(), error) {
	repoParts := strings.Split(actionID, "@")
	repoOrgAndName := repoParts[0]
	repoOrgAndNameSlice := strings.Split(repoOrgAndName, "/")
	if len(repoOrgAndNameSlice) != 2 {
		return "", func() {}, fmt.Errorf("actionID '%s' has invalid format", actionID)
	}
	repoName := repoOrgAndNameSlice[1]
	repoVersion := ""
	if len(repoParts) > 1 {
		repoVersion = repoParts[1]
	}

	repoURL := fmt.Sprintf("https://github.com/%s.git", repoOrgAndName)

	tempDir, err := pathutil.NormalizedOSTempDirPath("temp")
	if err != nil {
		return "", func() {}, err
	}

	cleanupFunc := func() {
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Println(err)
		}
	}

	cmd := exec.Command("git", "clone", repoURL)
	cmd.Dir = tempDir
	cmdLog, err := cmd.CombinedOutput()
	if err != nil {
		return "", cleanupFunc, fmt.Errorf("failed to clone action repository, error: %#v | output: %s", err, cmdLog)
	}

	repoDir := fmt.Sprintf("%s/%s", tempDir, repoName)

	if repoVersion != "" {
		cmd := exec.Command("git", "checkout", repoVersion)
		cmd.Dir = repoDir
		cmdLog, err := cmd.CombinedOutput()
		if err != nil {
			return "", cleanupFunc, fmt.Errorf("failed to checkout to the provided version of the action repository, error: %#v | output: %s", err, cmdLog)
		}
	}

	return repoDir, cleanupFunc, nil
}
