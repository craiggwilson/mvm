package internal

import (
	"archive/zip"
	"io"
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

	file, err := c.download(selected.Sha256, selected.URI, os.TempDir())
	if err != nil {
		return err
	}

	return c.decompressFile(file)
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

	c.writef("decompressing archive '%s' to '%s'", archive, target)

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
