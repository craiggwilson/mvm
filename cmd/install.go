package cmd

import (
	"fmt"
	"strings"

	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var installOpts = InstallOptions{
	RootOptions: &rootOpts,
}

func init() {
	installCmd.Flags().BoolVarP(&installOpts.Development, "development", "d", false, "include development versions")
	installCmd.Flags().BoolVarP(&installOpts.ReleaseCandidates, "releaseCandidates", "c", false, "include release candidates")

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a specific version of mongodb.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		installOpts.Version = args[0]
		return Install(installOpts)
	},
}

// InstallOptions are the options for installing a mongodb version.
type InstallOptions struct {
	*RootOptions

	Development       bool
	ReleaseCandidates bool

	Version string
}

// Select a mongodb version.
func Install(opts InstallOptions) error {
	versions, err := version.ListAll(opts.Config(), version.ListAllOptions{
		Development:       opts.Development,
		ReleaseCandidates: opts.ReleaseCandidates,
	})
	if err != nil {
		return err
	}

	var selected *version.Version
	for _, v := range versions {
		if v.Installed && strings.HasPrefix(v.Version.String(), opts.Version) {
			selected = v
			break
		}

		if selected == nil && strings.HasPrefix(v.Version.String(), opts.Version) {
			selected = v
		}
	}

	if selected == nil {
		return fmt.Errorf("no version matches %q", opts.Version)
	}

	return version.Install(opts.Config(), selected)
}
