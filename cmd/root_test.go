package cmd

import (
	"context"
	"os"
	"testing"

	"github.com/remydewolf/jumbler/pkg/config"
	"github.com/remydewolf/jumbler/pkg/jumbler"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {

	var actualConfig config.Config
	var actualMode jumbler.Mode

	originalCmd := runCmd
	runCmd = func(ctx context.Context, cfg config.Config, mode jumbler.Mode) error {
		actualConfig = cfg
		actualMode = mode
		return nil
	}
	pwdVal := "super_pwd_123"
	os.Setenv("JUMBLE_PWD", pwdVal)
	defer func() {
		runCmd = originalCmd
	}()
	defer os.Unsetenv("JUMBLE_PWD")

	testCases := []struct {
		name        string
		mode        jumbler.Mode
		dryRun      bool
		autoApprove bool
		quiet       bool
		password    string
		path        string
		args        []string
		err         string
	}{
		{
			name: "default",
			err:  "jumbler expects two parameters: [encrypt/decrypt] [path]",
			args: []string{""},
		},
		{
			name:        "--dry-run --auto-approve encrypt ./data/",
			mode:        jumbler.ModeEncrypt,
			dryRun:      true,
			autoApprove: true,
			quiet:       false,
			password:    pwdVal,
			path:        "./data/",
			args:        []string{"--dry-run", "--auto-approve", "encrypt", "./data/"},
		},
		{
			name:        "--quiet --password abc decrypt .",
			mode:        jumbler.ModeDecrypt,
			dryRun:      false,
			autoApprove: false,
			quiet:       true,
			password:    "abc",
			path:        ".",
			args:        []string{"--quiet", "--password", "abc", "decrypt", "."},
		},
		{
			name: "invalid .",
			err:  "invalid mode, expect encrypt or decrypt",
			args: []string{"invalid", "."},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := NewRootCmd()
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			if tc.err == "" {
				require.NoError(t, err)
				require.Equal(t, tc.mode, actualMode)
				require.Equal(t, tc.dryRun, actualConfig.DryRun)
				require.Equal(t, tc.autoApprove, actualConfig.AutoApprove)
				require.Equal(t, tc.password, actualConfig.Password)
				require.Equal(t, tc.path, actualConfig.Path)
				require.Equal(t, tc.quiet, actualConfig.Quiet)
			} else {
				require.ErrorContains(t, err, tc.err)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	cmd := NewRootCmd()
	cmd.SetArgs([]string{"--version"})
	require.NoError(t, cmd.Execute())
}
