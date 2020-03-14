package cmd

import (
	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var selectOpts = SelectOptions{
	RootOptions: rootOpts,
}

func init() {
	rootCmd.AddCommand(selectCmd)
}

var selectCmd = &cobra.Command{
	Use:   "select <version>",
	Short: "Select a specific version of mongodb.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		selectOpts.Version = args[0]
		return Select(selectOpts)
	},
}

// SelectOptions are the options for selecting a mongodb version.
type SelectOptions struct {
	RootOptions

	Version string
}

// Select a mongodb version.
func Select(opts SelectOptions) error {
	versions, err := version.ListInstalled(opts.Config())
	if err != nil {
		return err
	}

	matched, err := version.Match(versions, opts.Version)
	if err != nil {
		return err
	}

	return version.Select(opts.Config(), matched)
}
