package cmd

import (
	"github.com/micst/kxgctl/kxg"
	l "github.com/micst/kxgctl/kxg/logging"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a group address XML file from the config.",
	Long:  `Generate a group address XML file from the config.`,
	Run: func(cmd *cobra.Command, args []string) {
		kxg.LoadLibrary()
		kxg.LoadWorkspace()
		kxg.Validate()
		kxg.BuildTree()
		kxg.BuildAddresses()
		kxg.BuildDocument()
		if kxg.Args.XmlOut != "" {
			if kxg.Args.Dry {
				l.Info("dry run, not writing to disk")
			} else {
				kxg.Data.Document.WriteXml(kxg.Args.XmlOut)
			}
		} else {
			l.Info(kxg.Data.Document.GetXml())
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&kxg.Args.XmlOut, "out", "o", "", "generated xml group address file")
	generateCmd.Flags().StringVarP(&kxg.Args.ContextName, "context", "c", "", "context to use")
	generateCmd.Flags().BoolVar(&kxg.Args.SkipLibrary, "skip-library", false, "")
	generateCmd.Flags().BoolVar(&kxg.Args.Dry, "dry", false, "")
	generateCmd.Flags().BoolVar(&kxg.Args.SkipVerify, "skip-verify", false, "")
}
