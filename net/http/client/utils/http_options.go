package utils

import "net/http"

type HTTPOption func(req *http.Request)

func UserAgent(userAgent string) HTTPOption {
	return func(req *http.Request) {
		req.Header.Set("User-Agent", userAgent)
	}
}
