package utils

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func DownloadFile(ctx context.Context, url, destination string, opts ...HTTPOption) error {
	log.WithFields(log.Fields{
		"url":         url,
		"destination": destination,
	}).Trace("downloading file")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "error constructing HTTP Request object")
	}

	req.Header.Set("User-Agent", defaultUserAgent)

	for _, opt := range opts {
		opt(req)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error performing HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return errors.Errorf("unexpected status code: 200 or 206 expected while %d received", resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "downloaded_*.tmp")
	if err != nil {
		return errors.Wrap(err, "error creating temporary file")
	}
	log.WithFields(log.Fields{
		"path": tempFile.Name(),
	}).Trace("temporary file created")

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return errors.Wrap(err, "error writing to temporary file")
	}

	destinationDir := path.Dir(destination)
	log.WithFields(log.Fields{
		"path": destinationDir,
	}).Trace("ensuring directory structure")
	err = os.MkdirAll(destinationDir, 0o700)
	if err != nil {
		return errors.Wrap(err, "error creating directory structure for the final destination")
	}

	log.WithFields(log.Fields{
		"source":      tempFile.Name(),
		"destination": destination,
	}).Trace("moving the downloaded data")
	return errors.Wrap(
		os.Rename(tempFile.Name(), destination),
		"error moving temporary file to the final destination",
	)
}
