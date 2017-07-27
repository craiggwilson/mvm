package internal

import (
	"fmt"
	"os"
	"strings"
)

type InstallCmd struct {
	*RootCmd

	Version           string
	Development       bool
	ReleaseCandidates bool
}

func (c *InstallCmd) Execute() error {

	versions, err := c.allVersions(c.Development, c.ReleaseCandidates)
	if err != nil {
		return err
	}

	var selected *version
	for _, v := range versions {
		if v.Installed && strings.HasPrefix(v.Version.String(), c.Version) {
			selected = v
			break
		}

		if selected == nil && strings.HasPrefix(v.Version.String(), c.Version) {
			selected = v
		}
	}

	if selected == nil {
		return fmt.Errorf("no version matches '%s'", c.Version)
	}

	if selected.Installed {
		c.writef("not installing a version because '%s' matches '%s' and is already installed", selected.Version, c.Version)
		return nil
	}

	file, err := c.download(selected.Sha256, selected.URI, os.TempDir())
	if err != nil {
		return err
	}

	return c.decompressFile(file, c.VersionsPath)
}
