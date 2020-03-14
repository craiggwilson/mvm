package version

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/craiggwilson/mvm/pkg/config"
)

// Install the specified version.
func Install(cfg *config.Config, v *Version) error {
	if v.Installed {
		return fmt.Errorf("version %q is already installed", v.Version)
	}

	path, err := download(v.SHA256, v.URI, os.TempDir())
	if err != nil {
		return err
	}

	return decompressFile(path, cfg.InstallPath())
}

func decompressFile(archive, target string) error {
	log.Printf("[verbose] decompressing %q to %q\n", archive, target)
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}

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
