package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RunCmd struct {
	*RootCmd

	Version string
	Binary  string
	Args    []string
}

func (c *RunCmd) Execute() error {

	var err error
	var selected *version
	args := c.Args

	if len(args) > 0 {
		versions, err := c.installedVersions()
		if err != nil {
			return err
		}

		selected, err = c.selectVersion(versions, args[0])
		if err == nil {
			args = args[1:]
		}
	}

	if selected == nil {
		selected, err = c.activeVersion()
		if err != nil {
			return err
		}
	}

	switch c.Binary {
	case "mongod":
		port := "27017"
		portIdx := stringsIndex(args, "--port")
		if portIdx != -1 {
			if portIdx+1 > len(args) {
				return fmt.Errorf("--port has no value")
			}

			port = args[portIdx+1]
		}
		var dbpath string
		dbpathIdx := stringsIndex(args, "--dbpath")
		if dbpathIdx != -1 {
			if dbpathIdx+1 > len(args) {
				return fmt.Errorf("--dbpath has no value")
			}

			dbpath = args[dbpathIdx+1]
		} else {
			dbpath = filepath.Join(c.DataPath, strings.Replace(c.DataTemplate, "${port}", port, -1))
			args = append(args, "--dbpath", dbpath)
			c.verbosef("--dbpath did not exist, using '%s'", dbpath)
		}

		_, err := os.Stat(dbpath)
		if err != nil && os.IsNotExist(err) {
			c.verbosef("creating dbpath '%s'", dbpath)
			err = os.MkdirAll(dbpath, os.ModeDir)
		}
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(filepath.Join(selected.URI, c.Binary), args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
