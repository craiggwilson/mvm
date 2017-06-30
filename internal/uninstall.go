package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

type UninstallCmd struct {
	*RootCmd

	Version string
}

func (c *UninstallCmd) Execute() error {

	versions, err := c.installedVersions()
	if err != nil {
		return err
	}

	selected, err := c.selectVersion(versions, c.Version)
	if err != nil {
		return err
	}

	// require an exact match here
	if selected.Version.String() != c.Version {
		return fmt.Errorf("version '%s' is not installed", c.Version)
	}

	// selected.URI includes the bin folder... need to move one up
	parent := filepath.Dir(selected.URI)

	c.verbosef("removing '%s'", parent)

	return os.RemoveAll(parent)
}
