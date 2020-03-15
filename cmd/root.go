package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/craiggwilson/editline/pkg/editline"
	"github.com/craiggwilson/mvm/pkg/config"
	"github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var rootOpts RootOptions

func Execute(args []string) {
	rootCmd.PersistentFlags().BoolVar(&rootOpts.DisableColors, "nocolor", !isatty.IsTerminal(os.Stderr.Fd()), "disable colors forcefully")
	rootCmd.PersistentFlags().BoolVarP(&rootOpts.Verbose, "verbose", "v", false, "turn on verbose logging")

	rootCmd.SetArgs(args)
	_ = rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "mvm",
	Short: "MongoDB Version Manager",
	Long:  `MongoDB Version Manager`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		colors := rootOpts.Colors()
		editors := []editline.Editor{
			editline.Prefix("[info] ", editline.EditorFunc(func(line string) (string, editline.Action) {
				return line[7:], editline.ReplaceAction
			})),
			editline.Prefix("[verbose] ", editline.EditorFunc(func(line string) (string, editline.Action) {
				if !rootOpts.Verbose {
					return "", editline.RemoveAction
				}

				return colors.Cyan(line[10:]).String(), editline.ReplaceAction
			})),
		}

		log.SetFlags(0)
		log.SetOutput(editline.NewWriter(os.Stderr, editors...))
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		_ = log.Writer().(*editline.Writer).Flush()
	},
}

// RootOptions included in all other Options structs.
type RootOptions struct {
	DisableColors bool
	Verbose       bool
}

// Colors returns a colorizer for text.
func (o *RootOptions) Colors() aurora.Aurora {
	return aurora.NewAurora(!o.DisableColors)
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
