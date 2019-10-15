// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	GOBIN = ""

	CILINT_BINARY = "golangci-lint"
	CILINT_VERSION = "v1.19.1"
	CILINT_PATH = ""
	CILINT_REPOSITORY = "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

func init() {
	setupVars()
}

func setupVars() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	GOBIN = filepath.Join(cwd, "/bin")
	fmt.Printf("GOBIN: %s\n", GOBIN)
	if err := os.Setenv("GOBIN", GOBIN); err != nil {
		panic(err)
	}


	CILINT_PATH = filepath.Join(GOBIN, CILINT_BINARY)
}

// --------------------- Tools -----------------------

type Tools mg.Namespace

func (Tools) Install() {
	if !fileExists(CILINT_PATH) {
		cmd, args := toolInstallCmd(CILINT_REPOSITORY, CILINT_VERSION)
		if err := sh.Run(cmd, args...); err != nil {
			panic(err)
		}
	}
}

func toolInstallCmd(repository string, version string) (cmd string, args []string) {
	return "go", strings.Split(fmt.Sprintf("get -u %s@%s", repository, version), " ")
}

func fileExists(filepath string) bool {
	if runtime.GOOS == "windows" {
		filepath = fmt.Sprintf("%s.exe", filepath)
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}

	return true
}

// --------------------- Main Commands -----------------------

func Test() error {
	output, err := sh.Output("go", "test", "./...")
	fmt.Print(output)
	return err
}

func Lint() error {
	mg.Deps(Tools.Install)
	output, err := sh.Output(CILINT_PATH, "run", "./...")
	fmt.Print(output)
	return err
}

func CI() error {
	mg.Deps(Lint, Test)
	return nil
}