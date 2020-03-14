package version

import (
	"log"
	"os"
	"path/filepath"

	"github.com/craiggwilson/mvm/pkg/config"
)

// Uninstall the specified version.
func Uninstall(cfg *config.Config, v *Version) error {
	parent := filepath.Dir(v.URI)

	log.Printf("[info] removing %q\n", parent)

	return os.RemoveAll(parent)
}
