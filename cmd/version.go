package cmd

import (
	"time"

	"github.com/micst/kxgctl/kxg"
	l "github.com/micst/kxgctl/kxg/logging"

	"github.com/TwiN/go-color"
	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information.",
	Long:  `Show version information.`,
	Run: func(cmd *cobra.Command, args []string) {
		info := "kxgctl"
		if kxg.IsRelease && !versioninfo.DirtyBuild {
			info += " (version "
			info += kxg.Tag
			info += ")"
		} else {
			info += " ("
			info += color.InRed("development version")
			if versioninfo.DirtyBuild {
				info += ", "
				info += color.InRed("sources dirty!")
			}
			info += ")"
		}
		l.Info(info)
		l.Debug("commit: " + versioninfo.Revision)
		l.Debug("date:   " + versioninfo.LastCommit.Format(time.DateTime))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
