package version

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/craiggwilson/mvm/pkg/config"
)

func Installed(cfg *config.Config) ([]*Version, error) {
	log.Printf("[verbose] getting installed versions from %q\n", cfg.InstallPath())
	matches, err := filepath.Glob(filepath.Join(cfg.InstallPath(), "*-*.*.*"))
	if err != nil {
		return nil, err
	}

	activeVersion, err := Active(cfg)
	if err != nil && !errors.Is(err, ErrNoActiveVersion) {
		return nil, fmt.Errorf("failed getting active version: %w", err)
	}

	var versions []*Version
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			continue
		}

		v, err := FromPath(m)
		if err != nil {
			log.Printf("[info] skipping %q: %v\n", fi.Name(), err)
			continue
		}

		if v.URI == activeVersion.URI {
			v.Active = true
		}

		versions = append(versions, v)
	}

	Sort(versions)

	return versions, nil
}
