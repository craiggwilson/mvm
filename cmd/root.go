package cmd

import (
	"log"
	"os"

	"github.com/craiggwilson/mvm/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

var rootOpts RootOptions

func Execute(args []string) {
	rootCmd.SetArgs(args)
	_ = rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mvm",
	Short: "MongoDB Version Manager",
	Long:  `MongoDB Version Manager`,
}

type RootOptions struct{}

// Config returns the configuration for MVM.
func (o *RootOptions) Config() *config.Config {
	return &config.Config{
		Home: `C:\Users\craig\scoop\persist\mvm`,
	}
}
