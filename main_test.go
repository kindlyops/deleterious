package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestParseVersion(t *testing.T) {
	t.Parallel()

	cmd := exec.Command("/usr/bin/true")
	err := cmd.Run()
	if err != nil {
		t.Error("could not execute test", err)
	}
}
