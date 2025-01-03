package cmd

import (
	"github.com/spf13/cobra"

	"github.com/micst/kxgctl/kxg"

	l "github.com/micst/kxgctl/kxg/logging"
)

var useExample bool = false

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new kxg workspace.",
	Long:  `Initialize a new kxg workspace.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !useExample {
			kxg.LoadLibrary()
			if kxg.Lib.Exists {
				kxg.Ws.CopyFromDirectory(
					kxg.Lib.Directory,
					kxg.Args.Force,
				)
				return
			} else {
				l.Error("could not initialize workspace from library, falling back to example")
			}
		}
		if err := kxg.Ws.CopyFromResources(kxg.Args.Force); err != nil {
			l.Error("could not initialize workspace from resources into \"" + kxg.Ws.Directory)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&useExample, "example", false, "")
}
