package captaingithook

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
	"testing"
)

var expGitHooks = [...]string{
	"applypatch-msg",
	"pre-applypatch",
	"post-applypatch",
	"pre-commit",
	"prepare-commit-msg",
	"commit-msg",
	"post-commit",
	"pre-rebase",
	"post-checkout",
	"post-merge",
	"pre-push",
	"pre-receive",
	"update",
	"post-receive",
	"post-update",
	"push-to-checkout",
	"pre-auto-gc",
	"post-rewrite",
	"sendemail-validate",
}

const expHookFileScript = `#!/bin/sh
# captain-githook
# version ` + Version + `

hookName=` + "`basename \"$0\"`" + `
gitParams="$*"

if command -v captain-githook >/dev/null 2>&1; then
  captain-githook run --hook $hookName "$gitParams"
else
  echo "Can't find captain-githook, skipping $hookName hook"
  echo "You can reinstall it using 'go get -u github.com/swellaby/captain-githook' or delete this hook"
fi`

var expHookFileContents = []byte(expHookFileScript)

const gitHooksPath = "/usr/foo/repos/bar/.git/hooks"

func TestCreateAllHookFilesReturnsCorrectErrorOnHooksDirError(t *testing.T) {
	origGitFunc := getGitRepoHooksDirectory
	defer func() { getGitRepoHooksDirectory = origGitFunc }()
	getGitRepoHooksDirectory = func() (string, error) {
		return "", errors.New("")
	}

	assert.Equal(t, errInvalidGitHooksDirectoryPath, createAllHookFiles())
}

func TestCreateAllHookFilesReturnsCorrectErrorWhenSomeHooksNotCreated(t *testing.T) {
	expErrorHooks := [2]string{"pre-commit", "commit-msg"}
	expErrMsg := fmt.Sprintf("encountered an error while attempting to create one or more hook files. did not create hooks: %v", expErrorHooks)
	origGitFunc := getGitRepoHooksDirectory
	defer func() { getGitRepoHooksDirectory = origGitFunc }()
	getGitRepoHooksDirectory = func() (string, error) {
		return gitHooksPath, nil
	}
	originalWriteFile := writeFile
	defer func() { writeFile = originalWriteFile }()
	writeFile = func(filePath string, contents []byte) error {
		hook := strings.TrimPrefix(filePath, filepath.Join(gitHooksPath))
		hook = hook[1:len(hook)]

		if hook == "pre-commit" || hook == "commit-msg" {
			return errors.New("")
		}

		return nil
	}

	assert.Error(t, createAllHookFiles(), expErrMsg)
}

func TestCreateAllHookFilesCreatesCorrectHooks(t *testing.T) {
	var actHookPaths []string
	origGitFunc := getGitRepoHooksDirectory
	defer func() { getGitRepoHooksDirectory = origGitFunc }()
	getGitRepoHooksDirectory = func() (string, error) {
		return gitHooksPath, nil
	}
	originalWriteFile := writeFile
	defer func() { writeFile = originalWriteFile }()
	writeFile = func(filePath string, contents []byte) error {
		if !bytes.Equal(contents, expHookFileContents) {
			hook := strings.TrimPrefix(filePath, filepath.Join(gitHooksPath))
			hook = hook[1:len(hook)]
			t.Errorf("Incorrect script contents used for hook '%s'. Expected: %s, but got: %s", hook, string(expHookFileContents), string(contents))
		}
		actHookPaths = append(actHookPaths, filePath)
		return nil
	}
	createAllHookFiles()
	assert.Equal(t, len(actHookPaths), len(expGitHooks))
	for i, actHookPath := range actHookPaths {
		expHookPath := filepath.Join(gitHooksPath, expGitHooks[i])
		assert.Equal(t, expHookPath, actHookPath)
	}
}

func TestRemoveAllHookFilesReturnsCorrectError(t *testing.T) {
	var expErr error
	assert.Equal(t, expErr, removeAllHookFiles())
}
