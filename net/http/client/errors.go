// Package client implements HTTP utilities to ease multiple HTTP operations
//
// Deprecated: since first library version was implemented long ago Go ecosystem
// grew enough with more mature, tested and stable libraries. Please consider using
// another library for your purpose.
//
// In favor of backward compatibility this code won't be removed but it's frozen
// for adding new features. Only major issues will be fixed.
package client

import "github.com/pkg/errors"

var (
	ErrMisconfig            = errors.New("misconfiguration detected")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
)
