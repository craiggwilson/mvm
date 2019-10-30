package cmd

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/craiggwilson/mvm/pkg/version"
	"github.com/spf13/cobra"
)

var installOpts = InstallOptions{
	RootOptions: rootOpts,
}

func init() {
	installCmd.Flags().BoolVarP(&installOpts.Development, "development", "d", false, "include development versions")
	installCmd.Flags().BoolVarP(&installOpts.ReleaseCandidates, "releaseCandidates", "c", false, "include release candidates")

	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a specific version of mongodb.",
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		installOpts.Version = args[0]
		return Install(installOpts)
	},
}

// InstallOptions are the options for installing a mongodb version.
type InstallOptions struct {
	RootOptions

	Development       bool
	ReleaseCandidates bool

	Version string
}

// Select a mongodb version.
func Install(opts InstallOptions) error {
	versions, err := version.All(opts.Config(), version.RemoteOptions{
		Development:       opts.Development,
		ReleaseCandidates: opts.ReleaseCandidates,
	})
	if err != nil {
		return err
	}

	var selected *version.Version
	for _, v := range versions {
		if v.Installed && strings.HasPrefix(v.Version.String(), opts.Version) {
			selected = v
			break
		}

		if selected == nil && strings.HasPrefix(v.Version.String(), opts.Version) {
			selected = v
		}
	}

	if selected == nil {
		return fmt.Errorf("no version matches %q", opts.Version)
	}

	if selected.Installed {
		return fmt.Errorf("version %q is already installed", selected.Version)
	}

	path, err := version.Download(selected)
	if err != nil {
		return err
	}

	return decompressFile(path, opts.Config().InstallPath())
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
