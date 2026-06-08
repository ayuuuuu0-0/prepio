package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// New creates a reverse proxy handler for the given upstream base URL.
func New(target string) (http.Handler, error) {
	if len(target) == 0 {
		return nil, fmt.Errorf("proxy target is required")
	}

	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("parse proxy target: %w", err)
	}

	return httputil.NewSingleHostReverseProxy(u), nil
}
