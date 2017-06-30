package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/craiggwilson/mvm/internal"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("mvm", "A MongoDB version manager.")
	verbose = app.Flag("verbose", "write verbose output").Short('v').Bool()

	activePath   = app.Flag("activePath", "the symlink path for the active version").Hidden().Default(filepath.Join(os.Getenv("PROGRAMDATA"), "mvm", "active")).Envar(internal.MVMActiveEnvVarName).String()
	dataPath     = app.Flag("dataPath", "the path to store data").Hidden().Default(filepath.Join(os.Getenv("PROGRAMDATA"), "mvm", "data")).Envar(internal.MVMDataEnvVarName).String()
	dataTemplate = app.Flag("dataTemplate", "the data template for constructing a data directory").Hidden().Default(internal.MVMDataTemplateDefault).Envar(internal.MVMDataTemplateEnvVarName).String()
	versionsPath = app.Flag("versionsPath", "the path to store versions").Hidden().Default(filepath.Join(os.Getenv("PROGRAMDATA"), "mvm", "versions")).Envar(internal.MVMVersionsEnvVarName).String()

	env = app.Command("env", "lists the current environment as it pertains to MVM")

	list                 = app.Command("list", "list versions of mongodb")
	listAvailable        = list.Flag("available", "include the versions available").Short('a').Default("false").Bool()
	listDevelopment      = list.Flag("development", "include available development versions").Short('d').Default("false").Bool()
	listReleaseCandidate = list.Flag("releaseCandidates", "include available release candidates").Short('r').Default("false").Bool()

	run       = app.Command("run", "run mongodb binary")
	runBinary = run.Arg("binary", "the binary to run").Required().String()
	runArgs   = run.Arg("args", "remaining args").Strings()

	use        = app.Command("use", "use a specific version of mongodb")
	useVersion = use.Arg("version", "the version to use").Required().String()
)

type cmd interface {
	Execute() error
}

func main() {
	cmdName := kingpin.MustParse(app.Parse(os.Args[1:]))

	root := &internal.RootCmd{
		DataTemplate: *dataTemplate,
		ActivePath:   *activePath,
		DataPath:     *dataPath,
		VersionsPath: *versionsPath,
		Verbose:      *verbose,
		Writer:       os.Stdout,
	}

	err := root.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var cmd cmd
	switch cmdName {
	case env.FullCommand():
		cmd = &internal.EnvCmd{
			RootCmd: root,
		}
	case list.FullCommand():
		cmd = &internal.ListCmd{
			RootCmd:           root,
			Available:         *listAvailable,
			Development:       *listDevelopment,
			ReleaseCandidates: *listReleaseCandidate,
		}
	case run.FullCommand():
		cmd = &internal.RunCmd{
			RootCmd: root,
			Binary:  *runBinary,
			Args:    *runArgs,
		}
	case use.FullCommand():
		cmd = &internal.UseCmd{
			RootCmd: root,
			Version: *useVersion,
		}
	}

	err = cmd.Execute()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
