package cmd

import (
	"testing"
)

func TestRootCmd(t *testing.T) {
	// Root command should exist and have subcommands
	if rootCmd.Use != "ig" {
		t.Errorf("Use = %q, want %q", rootCmd.Use, "ig")
	}

	// Should have auth, media, comments, insights, discover, completion subcommands
	subCmds := rootCmd.Commands()
	names := make(map[string]bool)
	for _, cmd := range subCmds {
		names[cmd.Name()] = true
	}

	expected := []string{"auth", "media", "comments", "insights", "completion", "discover"}
	for _, name := range expected {
		if !names[name] {
			t.Errorf("missing subcommand %q", name)
		}
	}
}

func TestRootCmdFlags(t *testing.T) {
	f := rootCmd.PersistentFlags()

	if f.Lookup("account") == nil {
		t.Error("missing --account flag")
	}
	if f.Lookup("verbose") == nil {
		t.Error("missing --verbose flag")
	}
}

func TestPersistentPreRun(t *testing.T) {
	// Test with verbose off (default)
	verboseFlag = false
	rootCmd.PersistentPreRun(rootCmd, nil)

	// Test with verbose on
	verboseFlag = true
	rootCmd.PersistentPreRun(rootCmd, nil)

	// Reset
	verboseFlag = false
}

func TestExecuteFunc(t *testing.T) {
	// Test the Execute() wrapper for the success path (no args = prints help).
	// We can't test the error path since it calls os.Exit(1).
	rootCmd.SetArgs([]string{})
	Execute()
}
