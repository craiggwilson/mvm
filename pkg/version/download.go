package version

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"

	"github.com/cavaliercoder/grab"
)

func download(checksum, source, dest string) (string, error) {
	log.Printf("[verbose] downloading %q to %q\n", source, dest)
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

	log.Printf("[info] downloading %q\n", source)
	for {
		select {
		case <-t.C:
			log.Printf("[info] transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())
		case <-resp.Done:
			if err := resp.Err(); err != nil {
				return "", err
			}

			return resp.Filename, nil
		}
	}
}
