package main

import (
	"fmt"
	"os"

	"github.com/craiggwilson/mvm/internal"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app          = kingpin.New("mvm", "A MongoDB version manager.")
	verbose      = app.Flag("verbose", "write verbose output").Short('v').Bool()
	symlinkPath  = app.Flag("symlink-path", "the symlink path for the active version").Hidden().Envar(internal.MVMCurrentEnvVarName).String()
	mvmDirectory = app.Flag("mvm-directory", "the directory mvm places its resources").Hidden().Envar(internal.MVMEnvVarName).String()

	env = app.Command("env", "lists the current environment as it pertains to MVM")

	list    = app.Command("list", "list versions of mongodb")
	listAll = list.Flag("all", "list all the versions available").Default("false").Bool()

	use        = app.Command("use", "use a specific version of mongodb")
	useVersion = use.Arg("version", "the version to use").Required().String()
)

func main() {
	cmdName := kingpin.MustParse(app.Parse(os.Args[1:]))

	cfg := &internal.Config{
		SymlinkPath:  *symlinkPath,
		MVMDirectory: *mvmDirectory,
		Verbose:      *verbose,
		Writer:       os.Stdout,
	}

	err := cfg.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch cmdName {
	case env.FullCommand():
		err = internal.ExecuteEnv(&internal.EnvConfig{
			Config: cfg,
		})
	case list.FullCommand():
		err = internal.ExecuteList(&internal.ListConfig{
			Config: cfg,
			All:    *listAll,
		})
	case use.FullCommand():
		err = internal.ExecuteUse(&internal.UseConfig{
			Config:  cfg,
			Version: *useVersion,
		})
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
