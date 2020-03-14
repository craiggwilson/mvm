package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(dataCmd)
}

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Data provides information and actions about the data directories of the various MongoDB installs.",
}
