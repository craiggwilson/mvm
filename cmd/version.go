package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/craiggwilson/mvm/pkg/version"

	"github.com/spf13/cobra"
)

var binaryVersion string = "local-build"
var gitVersion string = "no-git-version"
var versionDate string = "no date"

var versionOpts = VersionOptions{
	RootOptions: &rootOpts,
	Out:         os.Stdout,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of the binary as well as the active version of mongodb",
	RunE: func(_ *cobra.Command, args []string) error {
		return Version(versionOpts)
	},
}

// VersionOptions are the options for selecting a mongodb version.
type VersionOptions struct {
	*RootOptions

	Out io.Writer
}

// Cleans a mongodb data directory.
func Version(opts VersionOptions) error {
	active, err := version.Active(opts.Config())

	var activeVersion string
	if err != nil && errors.Is(err, version.ErrNoActiveVersion) {
		activeVersion = "no active version"
	} else if err != nil {
		activeVersion = err.Error()
	} else {
		activeVersion = active.Version.String()
	}

	fmt.Fprintf(opts.Out, "{\n")
	fmt.Fprintf(opts.Out, "  \"version\"      : %q,\n", binaryVersion)
	fmt.Fprintf(opts.Out, "  \"git\"          : %q,\n", gitVersion)
	fmt.Fprintf(opts.Out, "  \"date\"         : %q,\n", versionDate)
	fmt.Fprintf(opts.Out, "  \"activeMongoDB\": %q\n", activeVersion)
	fmt.Fprintf(opts.Out, "}\n")

	return nil
}
