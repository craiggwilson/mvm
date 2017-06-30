package internal

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type InstallCmd struct {
	*RootCmd

	Version           string
	Development       bool
	ReleaseCandidates bool
}

func (c *InstallCmd) Execute() error {

	versions, err := c.allVersions(c.Development, c.ReleaseCandidates)
	if err != nil {
		return err
	}

	var selected *version
	for _, v := range versions {
		if strings.HasPrefix(v.Version.String(), c.Version) && v.Installed {
			selected = v
			break
		}

		if selected == nil && strings.HasPrefix(v.Version.String(), c.Version) {
			selected = v
		}
	}

	if selected.Installed {
		c.writef("not installing a version because '%s' matches '%s' and is already installed", selected.Version, c.Version)
		return nil
	}

	file, err := c.downloadFile(selected)
	if err != nil {
		return err
	}

	return c.decompressFile(file)
}

func (c *InstallCmd) downloadFile(selected *version) (string, error) {
	path := filepath.Join(os.TempDir(), selected.Sha256)

	if _, err := os.Stat(path); err == nil {
		c.verbosef("file for version '%s' already exists: '%s'", selected.Version, path)
		f, err := os.Open(path)
		if err != nil {
			return "", err
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return "", err
		}

		if hex.EncodeToString(h.Sum(nil)) == selected.Sha256 {
			return path, nil
		}

		c.verbosef("sha256 does not match, removing file '%s'", path)
		err = os.Remove(path)
		if err != nil {
			return "", err
		}
	}

	c.verbosef("downloading '%s' from '%s' and saving it to '%s'", selected.Version, selected.URI, path)
	resp, err := http.Get(selected.URI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := sha256.New()
	multi := io.MultiWriter(h, file)

	if _, err = io.Copy(multi, resp.Body); err != nil {
		return "", err
	}

	computed := hex.EncodeToString(h.Sum(nil))

	if computed != selected.Sha256 {
		return "", fmt.Errorf("sha256 did not match downloaded file: '%s' vs. '%s'", selected.Sha256, computed)
	}

	return path, nil
}

func (c *InstallCmd) decompressFile(archive string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	target := c.VersionsPath
	if err = os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}

	c.verbosef("decompressing archive '%s' to '%s'", archive, target)

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			if err = os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		pathDir := filepath.Dir(path)
		if err = os.MkdirAll(pathDir, os.ModePerm); err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return err
		}

		c.verbosef("decompressing file '%s'", file.Name)
		if _, err = io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			targetFile.Close()
			return err
		}

		fileReader.Close()
		targetFile.Close()
	}

	return nil
}
