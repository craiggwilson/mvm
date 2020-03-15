package cmd

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/craiggwilson/editline/pkg/editline"
	"github.com/craiggwilson/mvm/pkg/config"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var rootOpts RootOptions

func Execute(args []string) {
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.Verbose, "verbose", "v", false, "turn on verbose logging")

	rootCmd.SetArgs(args)
	_ = rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mvm",
	Short: "MongoDB Version Manager",
	Long:  `MongoDB Version Manager`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// don't log anything other than the message
		log.SetFlags(0)
		output := io.Writer(os.Stderr)

		verbose, _ := cmd.Flags().GetBool("verbose")

		editors := []editline.Editor{
			editline.Prefix("[info] ", editline.EditorFunc(func(line string) (string, editline.Action) {
				return aurora.Cyan(line[7:]).String(), editline.ReplaceAction
			})),
			editline.Prefix("[verbose] ", editline.EditorFunc(func(line string) (string, editline.Action) {
				if !verbose {
					return "", editline.RemoveAction
				}

				return aurora.Green(line[10:]).String(), editline.ReplaceAction
			})),
		}

		log.SetOutput(editline.NewWriter(output, editors...))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		w := log.Writer().(*editline.Writer)
		_ = w.Flush()
	},
}

type RootOptions struct {
	Verbose bool
}

// Config returns the configuration for MVM.
func (o *RootOptions) Config() *config.Config {
	home := os.Getenv("MVM")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			panic("user home directory not available")
		}

		home = filepath.Join(userHome, ".mvm")
	}

	return &config.Config{
		Home: home,
	}
}
