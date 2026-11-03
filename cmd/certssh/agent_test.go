package main

import (
	"bytes"
	"errors"
	"flag"
	"path/filepath"
	"strings"
	"testing"
)

func TestAgentInstallMarksManagedService(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CERTSSH_CONFIG_DIR", filepath.Join(home, ".config", "certssh"))
	t.Setenv("XDG_STATE_HOME", filepath.Join(home, ".local", "state"))

	var out, errOut bytes.Buffer
	if err := run([]string{"agent", "install", "-print"}, &out, &errOut); err != nil {
		t.Fatalf("agent install -print: %v\n%s", err, errOut.String())
	}
	if !strings.Contains(out.String(), managedServiceArg) {
		t.Fatalf("service definition does not mark the managed agent:\n%s", out.String())
	}
}

func TestAgentHelpHidesManagedService(t *testing.T) {
	for _, args := range [][]string{
		{"agent", "-h"},
		{"agent", managedServiceArg, "-h"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var out, errOut bytes.Buffer
			err := run(args, &out, &errOut)
			if !errors.Is(err, flag.ErrHelp) {
				t.Fatalf("run(%q) error = %v, want flag.ErrHelp", args, err)
			}
			if strings.Contains(errOut.String(), managedServiceArg) {
				t.Fatalf("agent help exposes %s:\n%s", managedServiceArg, errOut.String())
			}
		})
	}
}
