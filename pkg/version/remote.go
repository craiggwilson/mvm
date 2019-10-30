package version

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/blang/semver"
	"github.com/craiggwilson/mvm/pkg/config"
)

const (
	defaultEdition       = "enterprise"
	fullVersionsURI      = "https://downloads.mongodb.org/full.json"
	versionsTempFileName = "mvm_available_versions"
)

type RemoteOptions struct {
	OperatingSystem   string
	Edition           string
	Development       bool
	ReleaseCandidates bool
}

func Remote(cfg *config.Config, opts RemoteOptions) ([]*Version, error) {
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
