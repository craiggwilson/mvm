package version

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/blang/semver"
)

// ErrNoVersionsMatch is raised when no versions were matched.
var ErrNoVersionsMatch = errors.New("no versions match")

var versionRegex = regexp.MustCompile(`\d+\.\d+\.\d+.*`)

// FromPath returns the Version represented by the path.
func FromPath(path string) (*Version, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat failed on path %q: %w", path, err)
	}

	if fi.Name() == "bin" {
		return FromPath(filepath.Dir(path))
	}

	v := versionRegex.FindString(fi.Name())
	sv, err := semver.Parse(v)
	if err != nil {
		return nil, fmt.Errorf("semver parse failed on %q: %w", v, err)
	}

	return &Version{
		Version:   sv,
		URI:       filepath.Join(path, "bin"),
		Installed: true,
	}, nil
}

// Match a version given the fuzzy search string.
func Match(versions []*Version, fuzzy string) (*Version, error) {
	for _, v := range versions {
		if strings.HasPrefix(v.Version.String(), fuzzy) {
			return v, nil
		}
	}

	return nil, ErrNoVersionsMatch
}

// Version represents a version of MongoDB.
type Version struct {
	Version   semver.Version
	URI       string
	Installed bool
	Active    bool
	SHA256    string
}

// Sort sorts the versions in descending order.
func Sort(versions []*Version) {
	sort.Sort(versionSorter(versions))
}

type versionSorter []*Version

func (vs versionSorter) Len() int {
	return len(vs)
}

func (vs versionSorter) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs versionSorter) Less(i, j int) bool {
	return vs[i].Version.GT(vs[j].Version)
}
