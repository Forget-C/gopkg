package options

import (
	"github.com/spf13/pflag"
)

func NewPreInitFlags() *PreInitFlags {
	return &PreInitFlags{
		EnablePostScripts: true,
		EnablePreScripts:  true,
		BaseDir:           "./scripts",
		SkipError:         false,
		SkipRollbackError: true,
	}
}

type PreInitFlags struct {
	// BaseDir is the base directory for finding default files.
	BaseDir string
	// EnablePreRun is a flag to enable pre-run hooks.
	EnablePreScripts bool
	// EnablePostRun is a flag to enable post-run hooks.
	EnablePostScripts bool
	// SkipError is a flag to skip errors when exec error.
	SkipError bool
	// SkipRollbackError is a flag to skip errors when exec rollback error.
	SkipRollbackError bool
}

func (p *PreInitFlags) AddFlags(mainfs *pflag.FlagSet) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	defer func() {
		// Unhide deprecated flags. We want deprecated flags to show in Kubelet help.
		// We have some hidden flags, but we might as well unhide these when they are deprecated,
		// as silently deprecating and removing (even hidden) things is unkind to people who use them.
		fs.VisitAll(func(f *pflag.Flag) {
			if len(f.Deprecated) > 0 {
				f.Hidden = false
			}
		})
		mainfs.AddFlagSet(fs)
	}()
	fs.StringVar(&p.BaseDir, "base-dir", p.BaseDir, "The base directory for finding default files.")
	fs.BoolVar(&p.EnablePreScripts, "enable-pre-scripts", p.EnablePreScripts, "Enable pre-run scripts.")
	fs.BoolVar(&p.EnablePostScripts, "enable-post-scripts", p.EnablePostScripts, "Enable post-run scripts.")
	fs.BoolVar(&p.SkipError, "skip-error", p.SkipError, "Skip errors when exec error.")
	fs.BoolVar(&p.SkipRollbackError, "skip-rollback-error", p.SkipRollbackError, "Skip errors when exec rollback error.")
}
