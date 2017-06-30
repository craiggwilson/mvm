package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type version struct {
	Name      string
	Parts     []int
	URI       string
	Installed bool
	Active    bool
}

type versions []*version

func (vs versions) Len() int {
	return len(vs)
}

func (vs versions) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs versions) Less(i, j int) bool {
	a := vs[i].Parts
	b := vs[j].Parts

	for x := range a {
		if x == len(b) {
			return true
		}
		if b[x] < a[x] {
			return true
		}
	}
	return false
}

func activeVersion(cfg *Config) (*version, error) {
	versions, err := installedVersions(cfg)
	if err != nil {
		return nil, err
	}

	for _, v := range versions {
		if v.Active {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no version has been activated")
}

func installedVersions(cfg *Config) (versions, error) {
	dir := cfg.versionsDir()
	verbosef(cfg, "getting installed versions from %q", dir)

	matches, err := filepath.Glob(filepath.Join(dir, "*-?.?.?"))
	if err != nil {
		return nil, err
	}

	rgx := regexp.MustCompile("\\d\\.\\d\\.\\d.*")

	activeVersionPath := cfg.activeVersionPath()
	activeFile, _ := os.Stat(activeVersionPath)

	var versions versions
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			continue
		}

		name := rgx.FindString(fi.Name())
		parts := strings.Split(name, ".")
		var nameParts []int
		for _, p := range parts {
			np, _ := strconv.Atoi(p)
			nameParts = append(nameParts, np)
		}

		fi, _ = os.Stat(filepath.Join(m, "bin"))
		versions = append(versions, &version{
			Name:      name,
			Parts:     nameParts,
			URI:       filepath.Join(m, "bin"),
			Installed: true,
			Active:    activeFile != nil && os.SameFile(activeFile, fi),
		})
	}

	sort.Sort(versions)
	return versions, nil
}

func selectVersion(cfg *Config, target string) (*version, error) {
	versions, err := installedVersions(cfg)
	if err != nil {
		return nil, err
	}

	var selected *version
	for _, v := range versions {
		if strings.HasPrefix(v.Name, target) {
			verbosef(cfg, "selected version '%s'", v.Name)
			selected = v
			break
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no installed versions match '%s'", target)
	}

	return selected, nil
}
