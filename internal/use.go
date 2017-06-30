package internal

import "os"

func ExecuteUse(cfg *UseConfig) error {

	selected, err := selectVersion(cfg.Config, cfg.Version)
	if err != nil {
		return err
	}

	// remove symlink if it exists
	_, err = os.Stat(cfg.SymlinkPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil || !os.IsNotExist(err) {
		verbosef(cfg.Config, "removing existing path '%s'", cfg.SymlinkPath)
		err = os.Remove(cfg.SymlinkPath)
		if err != nil {
			return err
		}
	}

	// create hard link
	verbosef(cfg.Config, "creating symlink at '%s' to '%s'", cfg.SymlinkPath, selected.URI)
	err = os.Symlink(selected.URI, cfg.SymlinkPath)
	if err != nil {
		return err
	}

	return nil
}

type UseConfig struct {
	*Config

	Version string
}
