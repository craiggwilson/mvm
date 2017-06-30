package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ExecuteRun(cfg *RunConfig) error {

	var err error
	var selected *version
	args := cfg.Args

	if len(args) > 0 {
		selected, err = selectVersion(cfg.Config, args[0])
		if err == nil {
			args = args[1:]
		}
	}

	if selected == nil {
		selected, err = activeVersion(cfg.Config)
		if err != nil {
			return err
		}
	}

	switch cfg.Binary {
	case "mongod":
		port := "27017"
		portIdx := findArg(args, "--port")
		if portIdx != -1 {
			if portIdx+1 > len(args) {
				return fmt.Errorf("--port has no value")
			}

			port = args[portIdx+1]
		}
		var dbpath string
		dbpathIdx := findArg(args, "--dbpath")
		if dbpathIdx != -1 {
			if dbpathIdx+1 > len(args) {
				return fmt.Errorf("--dbpath has no value")
			}

			dbpath = args[dbpathIdx+1]
		} else {
			dbpath = filepath.Join(cfg.dataDir(), strings.Replace(cfg.DataTemplate, "${port}", port, -1))
			args = append(args, "--dbpath", dbpath)
			verbosef(cfg.Config, "--dbpath did not exist, using '%s'", dbpath)
		}

		_, err := os.Stat(dbpath)
		if err != nil && os.IsNotExist(err) {
			verbosef(cfg.Config, "creating dbpath '%s'", dbpath)
			err = os.MkdirAll(dbpath, os.ModeDir)
		}
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(filepath.Join(selected.URI, cfg.Binary), args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

type RunConfig struct {
	*Config

	Version string
	Binary  string
	Args    []string
}

func findArg(args []string, target string) int {
	return stringsIndex(args, target)
}
