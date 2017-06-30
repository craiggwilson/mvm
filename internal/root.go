package internal

import (
	"io"
	"os"
	"path/filepath"
)

type RootCmd struct {
	ActivePath   string
	DataPath     string
	VersionsPath string

	DataTemplate string

	Verbose bool
	Writer  io.Writer
}

func (c *RootCmd) Validate() error {
	if c.ActivePath != "" {
		p, err := filepath.Abs(c.ActivePath)
		if err != nil {
			return err
		}

		c.ActivePath = p
	}

	if c.DataPath != "" {
		p, err := filepath.Abs(c.DataPath)
		if err != nil {
			return err
		}
		c.DataPath = p
	}

	if c.VersionsPath != "" {
		p, err := filepath.Abs(c.VersionsPath)
		if err != nil {
			return err
		}
		c.VersionsPath = p
	}

	return nil
}

func (c *RootCmd) evalActivePath() string {
	_, err := os.Stat(c.ActivePath)
	if err != nil {
		return ""
	}

	p, err := filepath.EvalSymlinks(c.ActivePath)
	if err != nil {
		return ""
	}

	return p
}
