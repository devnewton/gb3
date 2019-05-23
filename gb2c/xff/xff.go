package xff

import (
	"net/http"
	"net/url"
)

//BuildURLFromRequest build an URL with X-Forwarded-* headers support
func BuildURLFromRequest(r *http.Request) url.URL {
	return url.URL{
		Scheme: schemeFromRequest(r),
		Host:   hostFromRequest(r),
	}
}

func schemeFromRequest(r *http.Request) string {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if len(scheme) == 0 {
		scheme = "http"
	}
	return scheme
}

func hostFromRequest(r *http.Request) string {
	host := r.Header.Get("X-Forwarded-Host")
	if len(host) == 0 {
		host = r.Host
	}
	return host
}
