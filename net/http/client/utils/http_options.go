// Package utils implements HTTP utilities to ease multiple HTTP operations
//
// Deprecated: since first library version was implemented long ago Go ecosystem
// grew enough with more mature, tested and stable libraries. Please consider using
// another library for your purpose.
//
// In favor of backward compatibility this code won't be removed but it's frozen
// for adding new features. Only major issues will be fixed.
package utils

import "net/http"

type HTTPOption func(req *http.Request)

func UserAgent(userAgent string) HTTPOption {
	return func(req *http.Request) {
		req.Header.Set("User-Agent", userAgent)
	}
}
