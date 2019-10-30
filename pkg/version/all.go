package version

import "github.com/craiggwilson/mvm/pkg/config"

// All gets all the versions, either installed or remote.
func All(cfg *config.Config, opts RemoteOptions) ([]*Version, error) {
	installedVersions, err := Installed(cfg)
	if err != nil {
		return nil, err
	}
	remoteVersions, err := Remote(cfg, opts)
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
