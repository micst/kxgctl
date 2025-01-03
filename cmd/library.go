package cmd

import (
	"github.com/spf13/cobra"

	"github.com/micst/kxgctl/kxg"

	l "github.com/micst/kxgctl/kxg/logging"
)

var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Initialize a new library.",
	Long:  `Initialize a new library.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := kxg.Lib.CopyFromResources(
			kxg.Args.Force,
		); err != nil {
			l.Error("could not initialize library in \"" + kxg.Lib.Directory + "\"")
		}
	},
}

func init() {
	initCmd.AddCommand(libraryCmd)
}
