package client

import "github.com/pkg/errors"

var (
	ErrMisconfig            = errors.New("misconfiguration detected")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
)
