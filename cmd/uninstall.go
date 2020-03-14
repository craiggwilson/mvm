package cmd

import (
	"fmt"

	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var uninstallOpts = UninstallOptions{
	RootOptions: rootOpts,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a specific version of mongodb.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		uninstallOpts.Version = args[0]
		return Uninstall(uninstallOpts)
	},
}

// UninstallOptions are the options for uninstalling a mongodb version.
type UninstallOptions struct {
	RootOptions

	Version string
}

// Select a mongodb version.
func Uninstall(opts UninstallOptions) error {
	versions, err := version.ListInstalled(opts.Config())
	if err != nil {
		return err
	}

	matched, err := version.Match(versions, opts.Version)
	if err != nil {
		return err
	}

	// require an exact match here
	if matched.Version.String() != opts.Version {
		return fmt.Errorf("version '%s' is not installed", opts.Version)
	}

	return version.Uninstall(opts.Config(), matched)
}
