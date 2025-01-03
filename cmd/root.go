package cmd

import (
	"os"
	"path/filepath"

	"github.com/micst/kxgctl/kxg"

	l "github.com/micst/kxgctl/kxg/logging"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kxgctl",
	Short: "A cli tool to generate and maintain KNX group addresses.",
	Long: `A cli tool to generate and maintain KNX group addresses.

See https://github.com/micst/kxgctl for details.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	lib := getLibraryDefaultDir()
	ws := getWorkspaceDefaultDir()
	rootCmd.PersistentFlags().StringVarP(&kxg.Ws.Directory, "workspace", "w", ws, "directory where to read/store all files.")
	rootCmd.PersistentFlags().StringVarP(&kxg.Lib.Directory, "library", "l", lib, "library directory.")
	rootCmd.PersistentFlags().BoolVarP(&kxg.Args.Force, "force", "f", false, "force overwriting existing files.")
	rootCmd.PersistentFlags().IntVarP(&l.Verbosity, "verbose", "v", 0, "show more log messages.")
	rootCmd.Flag("verbose").NoOptDefVal = "1"
	rootCmd.MarkFlagDirname("workspace")
	rootCmd.MarkFlagDirname("library")
}

func getLibraryDefaultDir() string {
	dir := ""
	if value, exists := os.LookupEnv("KXGCTL_LIBRARY"); exists {
		dir = value
	} else {
		if home, err := os.UserHomeDir(); err == nil {
			dir = filepath.Join(home, ".kxgctl")
		} else {
			l.Error("no user home found")
		}
	}
	return dir
}

func getWorkspaceDefaultDir() string {
	dir := ""
	if value, exists := os.LookupEnv("KXGCTL_WORKSPACE"); exists {
		dir = value
	} else {
		dir = "."
	}
	return dir
}
