package internal

import (
	"os"
	"path/filepath"
	"strings"
)

type EnvCmd struct {
	*RootCmd
}

func (c *EnvCmd) Execute() error {

	c.writef("%s=%s", MVMEnvVarName, os.Getenv(MVMEnvVarName))
	c.writef("%s=%s", MVMActiveEnvVarName, c.ActivePath)
	c.writef("%s=%s", MVMDataEnvVarName, c.DataPath)
	c.writef("%s=%s", MVMDataTemplateEnvVarName, c.DataTemplate)
	c.writef("%s=%s", MVMVersionsEnvVarName, c.VersionsPath)

	component := c.ActivePath
	path := os.Getenv("Path")
	pathParts := strings.Split(path, string(os.PathListSeparator))

	found := false
	for _, p := range pathParts {
		p = filepath.Clean(p)
		if p == component {
			found = true
			break
		}
	}
	if !found {
		c.write("")
		c.writef("Add '%s' to your PATH in order to get versioned binaries.", component)
	}

	return nil
}
