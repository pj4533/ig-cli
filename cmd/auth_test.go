package cmd

import (
	"testing"
)

func TestAuthCmd(t *testing.T) {
	if authCmd.Use != "auth" {
		t.Errorf("Use = %q, want %q", authCmd.Use, "auth")
	}

	subCmds := authCmd.Commands()
	names := make(map[string]bool)
	for _, cmd := range subCmds {
		names[cmd.Name()] = true
	}

	expected := []string{"setup", "add", "list", "remove"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing auth subcommand %q", name)
		}
	}
}

func TestAuthRemoveRequiresArg(t *testing.T) {
	err := authRemoveCmd.Args(authRemoveCmd, []string{})
	if err == nil {
		t.Error("expected error when no args provided")
	}
}

func TestAuthRemoveAcceptsOneArg(t *testing.T) {
	err := authRemoveCmd.Args(authRemoveCmd, []string{"testuser"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
