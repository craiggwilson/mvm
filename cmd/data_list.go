package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var dataListOpts = DataListOptions{
	RootOptions: rootOpts,
	Out:         os.Stdout,
}

func init() {
	dataCmd.AddCommand(dataListCmd)
}

var dataListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the information about the data of the installed versions of MongoDB.",
	RunE: func(_ *cobra.Command, args []string) error {
		return DataList(dataListOpts)
	},
}

// CleanOptions are the options for cleaning a data directory.
type DataListOptions struct {
	RootOptions
	Out io.Writer
}

// DataList lists the information about the data of the installed versions of MongoDB.
func DataList(opts DataListOptions) error {
	path := opts.Config().DataPath("", "")

	versionDirs, err := ioutil.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed reading %q: %w", path, err)
	}

	var versions []versionData

	for _, versionDir := range versionDirs {
		if !versionDir.IsDir() {
			continue
		}

		versionPath := filepath.Join(path, versionDir.Name())
		portDirs, err := ioutil.ReadDir(versionPath)
		if err != nil {
			return fmt.Errorf("failed reading %q: %w", versionPath, err)
		}

		var ports []portData
		for _, portDir := range portDirs {
			if !portDir.IsDir() {
				continue
			}

			portPath := filepath.Join(versionPath, portDir.Name())
			port, err := strconv.Atoi(portDir.Name())
			if err != nil {
				log.Printf("[verbose] unexpected port in data directory, %q: %v\n", portPath, err)
				continue
			}

			ports = append(ports, portData{
				Port: port,
				Path: portPath,
			})
		}
		versions = append(versions, versionData{
			Version: versionDir.Name(),
			Ports:   ports,
		})
	}

	if len(versions) == 0 {
		return nil
	}

	bytes, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed writing data: %w", err)
	}

	_, err = opts.Out.Write(bytes)
	return err
}

type versionData struct {
	Version string
	Ports   []portData
}

type portData struct {
	Port int
	Path string
}
