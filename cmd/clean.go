package cmd

import (
	"log"
	"os"

	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var cleanOpts = CleanOptions{
	RootOptions: rootOpts,
}

func init() {
	cleanCmd.Flags().StringVar(&cleanOpts.Port, "port", "", "the port")

	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean <version>",
	Short: "Cleans the data directory for a particular version and optionally a particular port.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		cleanOpts.Version = args[0]
		return Clean(cleanOpts)
	},
}

// CleanOptions are the options for cleaning a data directory.
type CleanOptions struct {
	RootOptions

	Version string
	Port    string
}

// Cleans a mongodb data directory.
func Clean(opts CleanOptions) error {
	versions, err := version.Installed(opts.Config())
	if err != nil {
		return err
	}

	matched, err := version.Match(versions, opts.Version)
	if err != nil {
		return err
	}

	dataPath := opts.Config().DataPath(matched.Version.String(), opts.Port)

	log.Printf("[info] removing data at %q", dataPath)
	return os.RemoveAll(dataPath)
}
