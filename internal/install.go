package internal

import "strings"

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
		if strings.HasPrefix(v.Version.String(), c.Version) && v.Installed {
			selected = v
			break
		}

		if selected == nil && strings.HasPrefix(v.Version.String(), c.Version) {
			selected = v
		}
	}

	if selected.Installed {
		c.writef("not installing a version because '%s' matches '%s' and is already installed", selected.Version, c.Version)
		return nil
	}

	c.writef("installing version '%s' from '%s'", selected.Version, selected.URI)
	return nil
}
