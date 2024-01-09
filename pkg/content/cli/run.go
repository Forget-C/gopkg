package cli

import (
	"fmt"
	"github.com/Forget-C/gopkg/pkg/log"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

func Run(cmd *cobra.Command) int {
	if logsInitialized, err := run(cmd); err != nil {

		if !logsInitialized {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			log.Logger.Error("command failed", err)
		}
		return 1
	}
	return 0
}

func RunNoErrOutput(cmd *cobra.Command) error {
	_, err := run(cmd)
	return err
}

func run(cmd *cobra.Command) (logsInitialized bool, err error) {
	rand.Seed(time.Now().UnixNano())

	cmd.SetGlobalNormalizationFunc(WordSepNormalizeFunc)

	if !cmd.SilenceUsage {
		cmd.SilenceUsage = true
		cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
			// Re-enable usage printing.
			c.SilenceUsage = false
			return err
		})
	}

	cmd.SilenceErrors = true

	// Inject logs.InitLogs after command line parsing into one of the
	// PersistentPre* functions.
	switch {
	case cmd.PersistentPreRun != nil:
		pre := cmd.PersistentPreRun
		cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
			pre(cmd, args)
		}
	case cmd.PersistentPreRunE != nil:
		pre := cmd.PersistentPreRunE
		cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			return pre(cmd, args)
		}
	default:
		cmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		}
	}

	err = cmd.Execute()
	return
}
