package config

import (
	"os"
	"path/filepath"
)

// Config is the configuration of the system.
type Config struct {
	Home string `yaml:"home"`
}

// ActivePath returns the path to the active version.
func (c *Config) ActivePath() string {
	path := filepath.Join(c.Home, "active")
	_, err := os.Stat(path)
	if err != nil {
		return ""
	}

	return path
}

// DataPath returns the path to the data for a version running on the specified port.
// If port is empty, it is left off the path.
func (c *Config) DataPath(version string, port string) string {
	path := filepath.Join(c.Home, "data", version)
	if port == "" {
		return path
	}

	return filepath.Join(path, "db"+port)
}

// InstallPath is the path to the installed versions.
func (c *Config) InstallPath() string {
	return filepath.Join(c.Home, "versions")
}
