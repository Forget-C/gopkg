package app

import (
	"fmt"
	"github.com/Forget-C/gopkg/cmd/executor/app/options"
	"github.com/Forget-C/gopkg/pkg/content/cli"
	"github.com/Forget-C/gopkg/pkg/log"
	"github.com/Forget-C/gopkg/pkg/scripts"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

func NewPreInitCommand() *cobra.Command {
	flagSet := pflag.NewFlagSet("executor", pflag.ExitOnError)
	flagSet.SetNormalizeFunc(cli.WordSepNormalizeFunc)
	preInitFlags := options.NewPreInitFlags()
	cmd := &cobra.Command{
		Use:   "executor",
		Short: "executor is script run/management tools",
		Long: `The executor is script run/management tools.
For example: initialize the database, initialize the configuration file, etc.
The program will find the scripts directory in the base-dir directory and execute all scripts in the scripts directory.
The script is executed in the following order: scripts in the pre_init directory-> scripts in the init directory-> scripts in the post-init directory.
The script will be executed according to the number sequence before the file name, for example: 001-init.sh -> 002-init.sh -> 003-init.sh.
If there is no number before the script file name, it will be executed in the dictionary order of the file name, for example: a-init.sh -> b-init.sh -> c-init.sh.
Any script execution fails and the program will exit.
`,
		SilenceUsage: true,
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return fmt.Errorf("parse args: %w", err)
			}
			stat, err := os.Stat(preInitFlags.BaseDir)
			if err != nil {
				return fmt.Errorf("stat base dir: %w", err)
			}
			if !stat.IsDir() {
				return fmt.Errorf("base dir is not dir")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.ParseFlags(args); err != nil {
				return fmt.Errorf("parse args: %w", err)
			}
			cmds := flagSet.Args()
			if len(cmds) > 0 {
				return fmt.Errorf("unknown command %+s", cmds[0])
			}
			manager := scripts.NewScriptManager(preInitFlags.BaseDir, log.Logger, cmd.OutOrStdout())
			if preInitFlags.EnablePreScripts {
				err := manager.AddPreScripts("pre_scripts")
				if err != nil {
					return fmt.Errorf("add pre scripts: %w", err)
				}
			}
			if preInitFlags.EnablePostScripts {
				err := manager.AddPostScripts("post_scripts")
				if err != nil {
					return fmt.Errorf("add post scripts: %w", err)
				}
			}
			err := manager.AddRunScripts("scripts")
			if err != nil {
				return fmt.Errorf("add run scripts: %w", err)
			}
			if err := manager.RunPrepare(); err != nil {
				return fmt.Errorf("run prepare: %w", err)
			}
			if err := manager.Execute(cmd.Context()); err != nil {
				return fmt.Errorf("\n preinit run scripts error: %w", err)
			}
			return nil
		},
	}
	preInitFlags.AddFlags(flagSet)
	cmd.Flags().AddFlagSet(flagSet)

	return cmd
}
