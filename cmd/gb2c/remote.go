package main

import (
	"net/http"
)

var gb2cHttpClient *http.Client

func init() {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	gb2cHttpClient = &http.Client{
		Transport: transport,
	}
}

func remoteHost(r *http.Request) string {
	host := r.Header.Get("X-Forwarded-Host")
	if len(host) == 0 {
		host = r.Host
	}
	return host
}

func remoteProtocol(r *http.Request) string {
	protocol := r.Header.Get("X-Forwarded-Proto")
	if len(protocol) == 0 {
		protocol = "http"
	}
	return protocol
}
