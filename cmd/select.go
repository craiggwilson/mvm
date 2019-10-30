package cmd

import (
	"log"
	"os"

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
	versions, err := version.Installed(opts.Config())
	if err != nil {
		return err
	}

	matched, err := version.Match(versions, opts.Version)
	if err != nil {
		return err
	}

	activePath := opts.Config().ActivePath()
	_, err = os.Stat(activePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil || !os.IsNotExist(err) {
		log.Printf("[verbose] removing existing path %q\n", activePath)
		err = os.Remove(activePath)
		if err != nil {
			return err
		}
	}

	log.Printf("[verbose] creating symlink at %q to %q\n", activePath, matched.URI)
	err = os.Symlink(matched.URI, activePath)
	if err != nil {
		return err
	}

	return nil
}
