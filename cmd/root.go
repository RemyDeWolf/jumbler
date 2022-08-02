package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/remydewolf/jumbler/pkg/config"
	"github.com/remydewolf/jumbler/pkg/jumbler"
	"github.com/remydewolf/jumbler/pkg/version"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	jumbleEnvVarPwd = "JUMBLE_PWD"
)

var runCmd = jumbler.Run

type rootOptions struct {
	mode        jumbler.Mode
	config      string
	dryRun      bool
	autoApprove bool
	quiet       bool
	password    string
	pathArg     string
	showVersion bool
}

func (ro rootOptions) loadConfig() (config.Config, error) {
	var cfg config.Config
	var err error
	if ro.config != "" {
		cfg, err = config.ReadFile(ro.config)
	} else {
		cfg, err = config.GetDefault()
	}
	cfg.DryRun = ro.dryRun
	if err != nil {
		return cfg, err
	}
	cfg.AutoApprove = ro.autoApprove

	if ro.password == "" {
		if val, ok := os.LookupEnv(jumbleEnvVarPwd); ok {
			cfg.Password = val
		}
	} else {
		cfg.Password = ro.password
	}
	if cfg.Password == "" {
		//if password is not set, prompt it
		if cfg.Password, err = enterPassword(); err != nil {
			return cfg, err
		}
	}

	if len(ro.pathArg) != 0 {
		//use path argument is specified
		cfg.Path = ro.pathArg
	}
	cfg.Quiet = ro.quiet

	return cfg, nil
}

func enterPassword() (string, error) {
	var password string
	verb := "Enter"
	for {
		fmt.Printf("\n%v password: ", verb)
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Printf("\n")
		if err != nil {
			return "", err
		}
		if len(bytePassword) == 0 {
			fmt.Println("password can't be empty")
			continue
		}
		if password != "" {
			if password != string(bytePassword) {
				return "", fmt.Errorf("passwords don't match")
			}
			break
		}
		password = string(bytePassword)
		verb = "Confirm"
	}
	return password, nil
}

// NewRootCmd returns the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	ro := rootOptions{}
	var rootCmd = &cobra.Command{
		Use:          "jumbler [encrypt/decrypt] [path]",
		Short:        "Quickly encode or decode a large number of file names",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if ro.showVersion {
				fmt.Println(formatVersion())
				return nil
			} else if len(args) != 2 {
				return fmt.Errorf("jumbler expects two parameters: [encrypt/decrypt] [path]")
			}
			var err error
			ro.mode, err = jumbler.GetMode(args[0])
			ro.pathArg = args[1]
			if err != nil {
				return err
			}
			cfg, err := ro.loadConfig()
			if err != nil {
				return err
			}
			return runCmd(cmd.Context(), cfg, ro.mode)
		},
	}

	flags := rootCmd.Flags()

	flags.StringVarP(&ro.config, "config", "c", "", "Config file")
	flags.StringVarP(&ro.password, "password", "p", "", fmt.Sprintf("Password used for encryption - also set %v to use an env variable instead", jumbleEnvVarPwd))
	flags.BoolVar(&ro.dryRun, "dry-run", false, "Preview only which files would be changed, no change made")
	flags.BoolVar(&ro.autoApprove, "auto-approve", false, " Skips interactive approval of plan before renaming files")
	flags.BoolVarP(&ro.quiet, "quiet", "q", false, " Do not print each file that would be modified")
	flags.BoolVarP(&ro.showVersion, "version", "v", false, "Show version")

	rootCmd.Commands()
	return rootCmd
}

func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func formatVersion() string {
	v := version.Get()
	if len(v.GitCommit) >= 7 {
		return fmt.Sprintf("%s+g%s", v.Version, v.GitCommit[:7])
	}
	return v.Version
}
