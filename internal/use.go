package internal

import (
	"fmt"
	"os"
	"strings"
)

func ExecuteUse(cfg *UseConfig) error {
	versions, err := installedVersions(cfg.Config)
	if err != nil {
		return err
	}

	var selected *version
	for _, v := range versions {
		if strings.HasPrefix(v.Name, cfg.Version) {
			verbosef(cfg.Config, "using version '%s'", v.Name)
			selected = v
			break
		}
	}

	if selected == nil {
		return fmt.Errorf("no installed versions match '%s'", cfg.Version)
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
