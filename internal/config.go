package internal

import (
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	DataTemplate string
	SymlinkPath  string
	MVMDirectory string
	Verbose      bool
	Writer       io.Writer
}

func (c *Config) Validate() error {
	if c.MVMDirectory != "" {
		mvm, err := filepath.Abs(c.MVMDirectory)
		if err != nil {
			return err
		}
		c.MVMDirectory = mvm
	}

	if c.SymlinkPath != "" {
		p, err := filepath.Abs(c.SymlinkPath)
		if err != nil {
			return err
		}

		c.SymlinkPath = p
	}

	return nil
}

func (c *Config) activeVersionPath() string {
	_, err := os.Stat(c.SymlinkPath)
	if err != nil {
		return ""
	}

	md, err := filepath.EvalSymlinks(c.SymlinkPath)
	if err != nil {
		return ""
	}

	return md
}

func (c *Config) dataDir() string {
	return filepath.Join(c.MVMDirectory, "data")
}

func (c *Config) versionsDir() string {
	return filepath.Join(c.MVMDirectory, "versions")
}
