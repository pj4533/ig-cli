package cmd

import (
	"testing"
)

func TestMediaCmd(t *testing.T) {
	if mediaCmd.Use != "media" {
		t.Errorf("Use = %q, want %q", mediaCmd.Use, "media")
	}

	subCmds := mediaCmd.Commands()
	names := make(map[string]bool)
	for _, cmd := range subCmds {
		names[cmd.Name()] = true
	}

	expected := []string{"list", "insights"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing media subcommand %q", name)
		}
	}
}

func TestCommentsCmd(t *testing.T) {
	if commentsCmd.Use != "comments" {
		t.Errorf("Use = %q, want %q", commentsCmd.Use, "comments")
	}

	subCmds := commentsCmd.Commands()
	names := make(map[string]bool)
	for _, cmd := range subCmds {
		names[cmd.Name()] = true
	}

	expected := []string{"list", "replies"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing comments subcommand %q", name)
		}
	}
}

func TestInsightsCmd(t *testing.T) {
	if insightsCmd.Use != "insights" {
		t.Errorf("Use = %q, want %q", insightsCmd.Use, "insights")
	}

	subCmds := insightsCmd.Commands()
	names := make(map[string]bool)
	for _, cmd := range subCmds {
		names[cmd.Name()] = true
	}

	expected := []string{"account", "audience"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing insights subcommand %q", name)
		}
	}
}

func TestDiscoverCmd(t *testing.T) {
	if discoverCmd.Use != "discover <username>" {
		t.Errorf("Use = %q, want %q", discoverCmd.Use, "discover <username>")
	}

	err := discoverCmd.Args(discoverCmd, []string{})
	if err == nil {
		t.Error("discover should require an argument")
	}

	err = discoverCmd.Args(discoverCmd, []string{"testuser"})
	if err != nil {
		t.Errorf("discover should accept one argument: %v", err)
	}
}

func TestMediaInsightsRequiresArg(t *testing.T) {
	err := mediaInsightsCmd.Args(mediaInsightsCmd, []string{})
	if err == nil {
		t.Error("media insights should require an argument")
	}
}

func TestCommentsListRequiresArg(t *testing.T) {
	err := commentsListCmd.Args(commentsListCmd, []string{})
	if err == nil {
		t.Error("comments list should require an argument")
	}
}

func TestCommentsRepliesRequiresArg(t *testing.T) {
	err := commentsRepliesCmd.Args(commentsRepliesCmd, []string{})
	if err == nil {
		t.Error("comments replies should require an argument")
	}
}

func TestCompletionCmd(t *testing.T) {
	if completionCmd.Use != "completion [bash|zsh|fish|powershell]" {
		t.Errorf("Use = %q", completionCmd.Use)
	}

	validArgs := completionCmd.ValidArgs
	if len(validArgs) != 4 {
		t.Errorf("ValidArgs length = %d, want 4", len(validArgs))
	}
}

func TestMediaListFlags(t *testing.T) {
	f := mediaListCmd.Flags()
	if f.Lookup("limit") == nil {
		t.Error("missing --limit flag on media list")
	}
}

func TestCommentsListFlags(t *testing.T) {
	f := commentsListCmd.Flags()
	if f.Lookup("limit") == nil {
		t.Error("missing --limit flag on comments list")
	}
}

func TestCommentsRepliesFlags(t *testing.T) {
	f := commentsRepliesCmd.Flags()
	if f.Lookup("limit") == nil {
		t.Error("missing --limit flag on comments replies")
	}
}

func TestInsightsAccountFlags(t *testing.T) {
	f := insightsAccountCmd.Flags()
	if f.Lookup("period") == nil {
		t.Error("missing --period flag on insights account")
	}
}
