package internal

import (
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	MongoDirectory string
	MVMDirectory   string
	Verbose        bool
	Writer         io.Writer
}

func (c *Config) Validate() error {
	if c.MVMDirectory != "" {
		mvm, err := filepath.Abs(c.MVMDirectory)
		if err != nil {
			return err
		}
		c.MVMDirectory = mvm
	}

	if c.MongoDirectory != "" {
		md, err := filepath.Abs(c.MongoDirectory)
		if err != nil {
			return err
		}

		c.MongoDirectory = md
	}

	return nil
}

func (c *Config) versionsDir() string {
	return filepath.Join(c.MVMDirectory, "versions")
}

func (c *Config) dataDir() string {
	return filepath.Join(c.MVMDirectory, "data")
}

func (c *Config) currentVersionPath() string {
	_, err := os.Stat(c.MongoDirectory)
	if err != nil {
		return ""
	}

	md, err := filepath.EvalSymlinks(c.MongoDirectory)
	if err != nil {
		return ""
	}

	return md
}
