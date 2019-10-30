package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var runOpts = RunOptions{
	RootOptions: rootOpts,
	Out:         os.Stdout,
	Err:         os.Stderr,
	In:          os.Stdin,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:                "run <binary>",
	Short:              "Run a mongodb binary.",
	DisableFlagParsing: true,
	RunE: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("must specify a binary")
		}

		runOpts.Binary = args[0]
		runOpts.Args = args[1:]
		return Run(runOpts)
	},
}

// RunOptions are the options for running a mongodb binary.
type RunOptions struct {
	RootOptions

	Binary string
	Args   []string

	Out io.Writer
	Err io.Writer
	In  io.Reader
}

// Run a mongodb binary.
func Run(opts RunOptions) error {
	args := opts.Args
	var matched *version.Version
	if len(args) > 0 {
		versions, err := version.Installed(opts.Config())
		if err != nil {
			return err
		}

		matched, err = version.Match(versions, args[0])
		if err == nil {
			args = args[1:]
		}
	}

	if matched == nil {
		var err error
		matched, err = version.Active(opts.Config())
		if err != nil {
			return err
		}
	}

	switch opts.Binary {
	case "mongod":
		port := "27017"
		portIdx := stringSliceIndex(args, "--port")
		if portIdx != -1 {
			if portIdx+1 > len(args) {
				return fmt.Errorf("--port has no value")
			}

			port = args[portIdx+1]
		}
		var dbpath string
		dbpathIdx := stringSliceIndex(args, "--dbpath")
		if dbpathIdx != -1 {
			if dbpathIdx+1 > len(args) {
				return fmt.Errorf("--dbpath has no value")
			}

			dbpath = args[dbpathIdx+1]
		} else {
			dbpath = opts.Config().DataPath(matched.Version.String(), port)
			args = append(args, "--dbpath", dbpath)
			log.Printf("[verbose] --dbpath was not specified, using %q\n", dbpath)
		}

		_, err := os.Stat(dbpath)
		if err != nil && os.IsNotExist(err) {
			log.Printf("[verbose] creating dbpath %q\n", dbpath)
			err = os.MkdirAll(dbpath, os.ModeDir)
		}
		if err != nil {
			return err
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	cmd := exec.Command(filepath.Join(matched.URI, opts.Binary), args...)

	cmd.Stdout = opts.Out
	cmd.Stderr = opts.Err
	cmd.Stdin = opts.In

	argStrings := make([]string, len(cmd.Args)-1)
	for i := 1; i < len(cmd.Args); i++ {
		argStrings[i-1] = cmd.Args[i]
		if strings.Contains(argStrings[i-1], " ") {
			argStrings[i-1] = "\"" + argStrings[i-1] + "\""
		}
	}

	log.Printf("[info] executing %s %s\n", cmd.Path, strings.Join(argStrings, " "))

	err := cmd.Start()
	if err != nil {
		return err
	}

	gracefulExit := false

	go func() {
		<-done
		gracefulExit = true
		_ = cmd.Process.Signal(os.Interrupt)
	}()

	err = cmd.Wait()
	if gracefulExit {
		return nil
	}

	return err
}

func stringSliceIndex(values []string, target string) int {
	for i, v := range values {
		if v == target {
			return i
		}
	}

	return -1
}
