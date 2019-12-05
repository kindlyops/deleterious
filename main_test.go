package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestParseVersion(t *testing.T) {
	t.Parallel()
	var err error
	if err != nil {
		t.Error("could not execute test", err)
	}
}
