package internal

import "sort"

func ExecuteList(cfg *ListConfig) error {

	versions, err := installedVersions(cfg.Config)
	if err != nil {
		return err
	}
	if cfg.All {
		// append(versions, availableVersions...)
		sort.Sort(versions)
	}

	write(cfg.Config, "")
	for _, v := range versions {
		s := ""
		if v.Active {
			s += "o"
		} else {
			s += " "
		}

		if v.Installed {
			s += "+"
		} else {
			s += " "
		}

		s += " "
		s += v.Name
		s += " (" + v.URI + ")"

		write(cfg.Config, s)
	}

	return nil
}

type ListConfig struct {
	*Config

	All bool
}
