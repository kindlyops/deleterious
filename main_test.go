package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

var cli = flag.String("cli", "", "The CLI binary")

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// Make sure that the x_defs override is working with the bazel build
func TestDefaultVersion(t *testing.T) {
	t.Parallel()

	if version != "dev" {
		t.Error("Expected", "dev", "got", version)
	}
}

func TestCLIVersion(t *testing.T) {
	t.Parallel()

	path, err := bazel.Runfile(*cli)
	if err != nil {
		t.Fatalf("Could not find runfile %s: %q", *cli, err)
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Missing binary %v", path)
	}
	file, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("Invalid filename %v", path)
	}
	cmd := exec.Command(file, "--version")
	cmd.Stderr = os.Stderr
	res, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed running '%v': %v", path, err)
	}
	output := strings.TrimSpace(string(res))

	if output != "dev.bazel version" {
		t.Error("Expected", "dev.bazel version", "got", output)
	}
}
