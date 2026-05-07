package integration

import (
	"strings"
	"testing"
)

func TestExploreListTables(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explore")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "users") || !strings.Contains(stdout, "orders") {
		t.Errorf("expected users and orders:\n%s", stdout)
	}
}

func TestExplainTable(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explain", "users")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "users") {
		t.Errorf("expected users in explain output:\n%s", stdout)
	}
}

func TestExplainWithDepth(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explain", "users", "-d", "2")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "orders") {
		t.Errorf("expected FK relationship to orders:\n%s", stdout)
	}
}

func TestExplainVerbose(t *testing.T) {
	env := Setup(t)
	env.SeedDefaults()

	stdout, _, exitCode := env.RunSquix("explain", "users", "-v")
	if exitCode != 0 {
		t.Fatalf("exit %d", exitCode)
	}
	if !strings.Contains(stdout, "users") {
		t.Errorf("expected verbose explain:\n%s", stdout)
	}
}
