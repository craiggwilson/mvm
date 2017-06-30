package internal

import (
	"os"
	"path/filepath"
	"strings"
)

func ExecuteEnv(cfg *EnvConfig) error {

	writef(cfg.Config, "set %s=%s", MVMEnvVarName, cfg.MVMDirectory)
	writef(cfg.Config, "set %s=%s", MVMActiveEnvVarName, cfg.SymlinkPath)

	component := filepath.Join(cfg.SymlinkPath)
	path := os.Getenv("Path")
	pathParts := strings.Split(path, string(os.PathListSeparator))

	found := false
	for _, p := range pathParts {
		p = filepath.Clean(p)
		if p == component {
			found = true
			break
		}
	}
	if !found {
		writef(cfg.Config, "path is missing '%s'", component)
	}

	return nil
}

type EnvConfig struct {
	*Config
}
