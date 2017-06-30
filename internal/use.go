package internal

import "os"

type UseCmd struct {
	*RootCmd

	Version string
}

func (c *UseCmd) Execute() error {

	selected, err := c.selectVersion(c.Version)
	if err != nil {
		return err
	}

	_, err = os.Stat(c.ActivePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil || !os.IsNotExist(err) {
		c.verbosef("removing existing path '%s'", c.ActivePath)
		err = os.Remove(c.ActivePath)
		if err != nil {
			return err
		}
	}

	c.verbosef("creating symlink at '%s' to '%s'", c.ActivePath, selected.URI)
	err = os.Symlink(selected.URI, c.ActivePath)
	if err != nil {
		return err
	}

	return nil
}
