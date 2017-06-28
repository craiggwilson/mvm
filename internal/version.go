package internal

import (
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
	Current   bool
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

func installedVersions(cfg *Config) (versions, error) {
	dir := cfg.versionsDir()
	writef(cfg, "getting installed versions from %q", dir)

	matches, err := filepath.Glob(filepath.Join(dir, "*-?.?.?"))
	if err != nil {
		return nil, err
	}

	rgx := regexp.MustCompile("\\d\\.\\d\\.\\d.*")

	currentVersionPath := cfg.currentVersionPath()

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
		versions = append(versions, &version{
			Name:      name,
			Parts:     nameParts,
			URI:       m,
			Installed: true,
			Current:   m == currentVersionPath,
		})
	}

	sort.Sort(versions)
	return versions, nil
}
