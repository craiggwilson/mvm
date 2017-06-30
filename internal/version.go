package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/blang/semver"
)

type version struct {
	Version   semver.Version
	URI       string
	Installed bool
	Active    bool
	Sha256    string
}

type versions []*version

func (vs versions) Len() int {
	return len(vs)
}

func (vs versions) Swap(i, j int) {
	vs[i], vs[j] = vs[j], vs[i]
}

func (vs versions) Less(i, j int) bool {
	return vs[i].Version.GT(vs[j].Version)
}

func (c *RootCmd) activeVersion() (*version, error) {
	versions, err := c.installedVersions()
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

func (c *RootCmd) allVersions(dev, rc bool) (versions, error) {
	installedVersions, err := c.installedVersions()
	if err != nil {
		return nil, err
	}
	availableVersions, err := c.availableVersions(dev, rc)
	if err != nil {
		return nil, err
	}

	var versions versions
	i := 0
	j := 0
	for {
		if i < len(installedVersions) && j < len(availableVersions) {
			iv, jv := installedVersions[i], availableVersions[j]
			if iv.Version.EQ(jv.Version) {
				versions = append(versions, iv)
				i++
				j++
			} else if iv.Version.GT(jv.Version) {
				versions = append(versions, iv)
				i++
			} else {
				versions = append(versions, jv)
				j++
			}

			continue
		}

		if i < len(installedVersions) {
			versions = append(versions, installedVersions[i:]...)
		} else if j < len(availableVersions) {
			versions = append(versions, availableVersions[j:]...)
		}

		break
	}

	return versions, nil
}

func (c *RootCmd) availableVersions(dev bool, rc bool) (versions, error) {
	c.verbosef("hitting '%s' for available versions", fullVersionsURI)

	client := http.DefaultClient
	resp, err := client.Get(fullVersionsURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	root := &struct {
		Versions []struct {
			DevelopmentRelease bool `json:"development_release"`
			Downloads          []struct {
				Archive struct {
					Sha256 string `json:"sha256"`
					URL    string `json:"url"`
				}
				Edition string `json:"edition"`
				Target  string `json:"target"`
			} `json:"downloads"`
			ProductionRelease bool   `json:"production_release"`
			ReleaseCandidate  bool   `json:"release_candidate"`
			Version           string `json:"version"`
		} `json:"versions"`
	}{}

	err = json.Unmarshal(body, root)
	if err != nil {
		return nil, err
	}

	var versions versions
	for _, v := range root.Versions {

		if !v.ProductionRelease && !(dev && v.DevelopmentRelease) && !(rc && v.ReleaseCandidate) {
			continue
		}

		sv, err := semver.Parse(v.Version)
		if err != nil {
			c.verbosef("skipping '%s': %v", v.Version, err)
			continue
		}

		for _, d := range v.Downloads {
			if d.Target != versionTarget {
				continue
			}

			if d.Edition != "enterprise" {
				continue
			}

			versions = append(versions, &version{
				Version:   sv,
				URI:       d.Archive.URL,
				Installed: false,
				Active:    false,
				Sha256:    d.Archive.Sha256,
			})
		}
	}

	sort.Sort(versions)

	return versions, nil
}

func (c *RootCmd) installedVersions() (versions, error) {
	c.verbosef("getting installed versions from '%s'", c.VersionsPath)

	matches, err := filepath.Glob(filepath.Join(c.VersionsPath, "*-?.?.?"))
	if err != nil {
		return nil, err
	}

	rgx := regexp.MustCompile("\\d\\.\\d\\.\\d.*")

	activePath := c.evalActivePath()
	activeFile, _ := os.Stat(activePath)

	var versions versions
	for _, m := range matches {
		fi, err := os.Stat(m)
		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			continue
		}

		v := rgx.FindString(fi.Name())
		sv, err := semver.Parse(v)
		if err != nil {
			c.verbosef("skipping '%s': %v", v, err)
			continue
		}
		fi, _ = os.Stat(filepath.Join(m, "bin"))
		versions = append(versions, &version{
			Version:   sv,
			URI:       filepath.Join(m, "bin"),
			Installed: true,
			Active:    activeFile != nil && os.SameFile(activeFile, fi),
		})
	}

	sort.Sort(versions)
	return versions, nil
}

func (c *RootCmd) selectVersion(target string) (*version, error) {
	versions, err := c.installedVersions()
	if err != nil {
		return nil, err
	}

	var selected *version
	for _, v := range versions {
		if strings.HasPrefix(v.Version.String(), target) {
			c.verbosef("selected version '%s'", v.Version)
			selected = v
			break
		}
	}

	if selected == nil {
		return nil, fmt.Errorf("no installed versions match '%s'", target)
	}

	return selected, nil
}
