package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/craiggwilson/editline/pkg/editline"
	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var listOpts = ListOptions{
	RootOptions: &rootOpts,
	Out:         os.Stdout,
}

func init() {
	listCmd.Flags().BoolVarP(&listOpts.Remote, "remote", "r", false, "include remote versions that haven't been downloaded")
	listCmd.Flags().BoolVarP(&listOpts.Development, "development", "d", false, "include development versions")
	listCmd.Flags().BoolVarP(&listOpts.ReleaseCandidates, "releaseCandidates", "c", false, "include release candidates")

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the versions of mongodb available.",
	PreRun: func(_ *cobra.Command, args []string) {
		colors := listOpts.RootOptions.Colors()
		listOpts.Out = editline.NewWriter(listOpts.Out,
			editline.Prefix("o +", editline.EditorFunc(func(line string) (string, editline.Action) {
				return colors.BrightGreen(line).Bold().String(), editline.ReplaceAction
			})),
			editline.Prefix("   ", editline.EditorFunc(func(line string) (string, editline.Action) {
				return colors.Yellow(line).Bold().String(), editline.ReplaceAction
			})),
		)
	},
	RunE: func(_ *cobra.Command, args []string) error {
		return List(listOpts)
	},
	PostRun: func(_ *cobra.Command, args []string) {
		_ = listOpts.Out.(*editline.Writer).Flush()
	},
}

// ListOptions are the options for listing versions.
type ListOptions struct {
	*RootOptions

	Remote            bool
	Development       bool
	ReleaseCandidates bool

	Out io.Writer
}

// List the versions.
func List(opts ListOptions) error {
	var versions []*version.Version
	var err error
	if opts.Remote {
		versions, err = version.ListAll(opts.Config(), version.ListAllOptions{
			Development:       opts.Development,
			ReleaseCandidates: opts.ReleaseCandidates,
		})
	} else {
		versions, err = version.ListInstalled(opts.Config())
	}
	if err != nil {
		return err
	}

	for _, v := range versions {
		activeFlag := " "
		if v.Active {
			activeFlag = "o"
		}
		installedFlag := " "
		if v.Installed {
			installedFlag = "+"
		}
		path := ""
		if opts.Verbose {
			path = v.URI
		}

		fmt.Fprintf(opts.Out, "%s %s %-12s %s\n", activeFlag, installedFlag, v.Version.String(), path)
	}

	return nil
}
