package version

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/craiggwilson/mvm/pkg/config"
)

// ErrNoActiveVersion is a sentinel error returned when no active version exists.
var ErrNoActiveVersion = errors.New("no active version")

// Active returns the currently active version.
func Active(cfg *config.Config) (*Version, error) {
	path, err := filepath.EvalSymlinks(cfg.ActivePath())
	if err != nil {
		return nil, err
	}

	v, err := FromPath(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoActiveVersion, err)
	}

	v.Active = true
	return v, nil
}
