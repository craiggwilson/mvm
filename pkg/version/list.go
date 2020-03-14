package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/blang/semver"
	"github.com/craiggwilson/mvm/pkg/config"
)

type ListAllOptions struct {
	OperatingSystem   string
	Edition           string
	Development       bool
	ReleaseCandidates bool
}

// ListAll gets all the versions, either installed or remote.
func ListAll(cfg *config.Config, opts ListAllOptions) ([]*Version, error) {
	installedVersions, err := ListInstalled(cfg)
	if err != nil {
		return nil, err
	}
	remoteVersions, err := ListRemote(cfg, ListRemoteOptions{
		OperatingSystem:   opts.OperatingSystem,
		Edition:           opts.Edition,
		Development:       opts.Development,
		ReleaseCandidates: opts.ReleaseCandidates,
	})

	if err != nil {
		return nil, err
	}

	var versions []*Version
	i := 0
	j := 0
	for {
		if i < len(installedVersions) && j < len(remoteVersions) {
			iv, jv := installedVersions[i], remoteVersions[j]
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
		} else if j < len(remoteVersions) {
			versions = append(versions, remoteVersions[j:]...)
		}

		break
	}

	return versions, nil
}

func ListInstalled(cfg *config.Config) ([]*Version, error) {
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

const (
	defaultEdition       = "enterprise"
	fullVersionsURI      = "https://downloads.mongodb.org/full.json"
	versionsTempFileName = "mvm_available_versions"
)

type ListRemoteOptions struct {
	OperatingSystem   string
	Edition           string
	Development       bool
	ReleaseCandidates bool
}

func ListRemote(cfg *config.Config, opts ListRemoteOptions) ([]*Version, error) {
	if opts.Edition == "" {
		opts.Edition = defaultEdition
	}
	if opts.OperatingSystem == "" {
		opts.OperatingSystem = runtime.GOOS
	}

	versionFile := filepath.Join(os.TempDir(), versionsTempFileName)

	needNewFile := true
	fi, err := os.Stat(versionFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		if fi.ModTime().Add(24 * time.Hour).Before(time.Now()) {
			log.Printf("[verbose] found cached available versions file, but it is older than 24 hours\n")
			_ = os.Remove(versionFile)
		} else {
			needNewFile = false
		}
	}

	if needNewFile {
		_, err = download("", fullVersionsURI, versionFile)
		if err != nil {
			return nil, err
		}
	}

	body, err := ioutil.ReadFile(versionFile)
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

	var versions []*Version
	for _, v := range root.Versions {

		if !v.ProductionRelease && !(opts.Development && v.DevelopmentRelease) && !(opts.ReleaseCandidates && v.ReleaseCandidate) {
			continue
		}

		sv, err := semver.Parse(v.Version)
		if err != nil {
			log.Printf("[verbose] skipping %q: %v\n", v.Version, err)
			continue
		}

		for _, d := range v.Downloads {
			if d.Target != opts.OperatingSystem {
				continue
			}

			if d.Edition != opts.Edition {
				continue
			}

			versions = append(versions, &Version{
				Version: sv,
				URI:     d.Archive.URL,
				SHA256:  d.Archive.Sha256,
			})
		}
	}

	Sort(versions)

	return versions, nil
}
