package version

import (
	"log"
	"os"

	"github.com/craiggwilson/mvm/pkg/config"
)

// Select the specified version.
func Select(cfg *config.Config, v *Version) error {
	activePath := cfg.ActivePath()
	_, err := os.Stat(activePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil || !os.IsNotExist(err) {
		log.Printf("[verbose] removing existing path %q\n", activePath)
		err = os.Remove(activePath)
		if err != nil {
			return err
		}
	}

	log.Printf("[verbose] creating symlink at %q to %q\n", activePath, v.URI)
	err = os.Symlink(v.URI, activePath)
	if err != nil {
		return err
	}

	return nil
}
