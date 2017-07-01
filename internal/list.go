package internal

type ListCmd struct {
	*RootCmd

	Development       bool
	ReleaseCandidates bool
	Remote            bool
}

func (c *ListCmd) Execute() error {

	var versions versions
	var err error
	if c.Remote {
		versions, err = c.allVersions(c.Development, c.ReleaseCandidates)
	} else {
		versions, err = c.installedVersions()
	}

	if err != nil {
		return err
	}

	c.write("")
	for _, v := range versions {
		s := ""
		if v.Active {
			s += "o"
		} else {
			s += " "
		}

		s += " "

		if v.Installed {
			s += "+"
		} else {
			s += " "
		}

		s += " "
		s += v.Version.String()
		s += " (" + v.URI + ")"

		c.write(s)
	}

	return nil
}
