package cmd

import (
	"github.com/micst/kxgctl/kxg"

	"github.com/spf13/cobra"
)

var Show = struct {
	Contexts   bool
	Attributes bool
	Templates  bool
	Devices    bool
}{
	Contexts:   false,
	Attributes: false,
	Templates:  false,
	Devices:    false,
}

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Show certain aspects of kxgctl configuration.",
	Long:  `Show certain aspects of kxgctl configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		kxg.LoadLibrary()
		kxg.LoadWorkspace()
		kxg.Validate()
		if Show.Contexts {
			kxg.ShowContexts()
		}
		if Show.Attributes {
			kxg.ShowAttributes()
		}
		if Show.Templates {
			kxg.ShowTemplates()
		}
		if Show.Devices {
			kxg.ShowDevices()
		}
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().BoolVar(&kxg.Args.SkipLibrary, "skip-library", false, "")
	inspectCmd.Flags().BoolVar(&kxg.Args.SkipVerify, "skip-verify", false, "")
	inspectCmd.Flags().BoolVarP(&Show.Contexts, "contexts", "c", false, "")
	inspectCmd.Flags().BoolVarP(&Show.Attributes, "attributes", "a", false, "")
	inspectCmd.Flags().BoolVarP(&Show.Templates, "templates", "t", false, "")
	inspectCmd.Flags().BoolVarP(&Show.Devices, "devices", "d", false, "")
	inspectCmd.Flags().StringVar(&kxg.Args.ContextName, "context", "", "context to use")
}
