package internal

import (
	"archive/zip"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"time"

	"encoding/hex"

	"github.com/cavaliercoder/grab"
)

type RootCmd struct {
	ActivePath   string
	DataPath     string
	VersionsPath string

	DataTemplate string

	Verbose bool
	Writer  io.Writer
}

func (c *RootCmd) Validate() error {
	if c.ActivePath != "" {
		p, err := filepath.Abs(c.ActivePath)
		if err != nil {
			return err
		}

		c.ActivePath = p
	}

	if c.DataPath != "" {
		p, err := filepath.Abs(c.DataPath)
		if err != nil {
			return err
		}
		c.DataPath = p
	}

	if c.VersionsPath != "" {
		p, err := filepath.Abs(c.VersionsPath)
		if err != nil {
			return err
		}
		c.VersionsPath = p
	}

	return nil
}

func (c *InstallCmd) decompressFile(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

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

func (c *RootCmd) download(checksum, source, dest string) (string, error) {
	client := grab.NewClient()
	req, err := grab.NewRequest(dest, source)
	if err != nil {
		return "", err
	}

	if checksum != "" {
		cs, err := hex.DecodeString(checksum)
		if err != nil {
			return "", err
		}
		req.SetChecksum(sha256.New(), cs, true)
	}

	resp := client.Do(req)

	t := time.NewTicker(time.Second)
	defer t.Stop()

	c.writef("downloading '%s'", source)
	for {
		select {
		case <-t.C:
			c.writef("  %f%% complete", resp.Progress()*100)
		case <-resp.Done:
			if err := resp.Err(); err != nil {
				return "", err
			}

			return resp.Filename, nil
		}
	}
}

func (c *RootCmd) evalActivePath() string {
	_, err := os.Stat(c.ActivePath)
	if err != nil {
		return ""
	}

	p, err := filepath.EvalSymlinks(c.ActivePath)
	if err != nil {
		return ""
	}

	return p
}
