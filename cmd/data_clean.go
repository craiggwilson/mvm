package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var dataCleanOpts = DataCleanOptions{
	RootOptions: rootOpts,
}

func init() {
	dataCleanCmd.Flags().StringVar(&dataCleanOpts.Port, "port", "", "the port")

	dataCmd.AddCommand(dataCleanCmd)
}

var dataCleanCmd = &cobra.Command{
	Use:   "clean <version>",
	Short: "Cleans the data directory for a particular version and optionally a particular port.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		dataCleanOpts.Version = args[0]
		return Clean(dataCleanOpts)
	},
}

// DataCleanOptions are the options for cleaning a data directory.
type DataCleanOptions struct {
	RootOptions

	Version string
	Port    string
}

// Cleans a mongodb data directory.
func Clean(opts DataCleanOptions) error {
	versions, err := version.Installed(opts.Config())
	if err != nil {
		return err
	}

	matched, err := version.Match(versions, opts.Version)
	if err != nil {
		return err
	}

	if matched.Version.String() != opts.Version {
		return fmt.Errorf("data for version '%s' does not exist", opts.Version)
	}

	dataPath := opts.Config().DataPath(matched.Version.String(), opts.Port)

	log.Printf("[info] removing data at %q", dataPath)
	return os.RemoveAll(dataPath)
}
