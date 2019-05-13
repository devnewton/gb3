package linuxfr

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var clientID string
var clientSecret string

func init() {
	flag.StringVar(&clientID, "linuxfr-client-id", "", "linuxfr.org's API client id")
	flag.StringVar(&clientSecret, "linuxfr-client-secret", "", "linuxfr.org's API client secret")
}

//RegisterLinuxfrAPI setup http linuxfr API endpoints
func RegisterLinuxfrAPI() {
	http.HandleFunc("/api/linuxfr/authorize", func(w http.ResponseWriter, r *http.Request) {
		redirectURL := url.URL{
			Scheme: r.URL.Scheme,
			Host:   r.URL.Host,
			Path:   "/api/linuxfr/connected",
		}

		linuxfrAuthorizeURL := url.URL{
			Scheme: "https",
			Host:   "linuxfr.org",
			Path:   "/api/oauth/authorize",
		}
		linuxfrAuthorizeQuery := linuxfrAuthorizeURL.Query()
		linuxfrAuthorizeQuery.Set("client_id", clientID)
		linuxfrAuthorizeQuery.Set("response_type", "code")
		linuxfrAuthorizeQuery.Set("scope", "account+board")
		linuxfrAuthorizeQuery.Set("redirect_uri", redirectURL.String())

		http.Redirect(w, r, linuxfrAuthorizeURL.String(), 303)
	})
	http.HandleFunc("/api/linuxfr/connected", func(w http.ResponseWriter, r *http.Request) {
		redirectURL := url.URL{
			Scheme: r.URL.Scheme,
			Host:   r.URL.Host,
			Path:   "/linuxfr/connected.html",
		}
		tokenParams := url.Values{}
		tokenParams.Set("client_id", clientID)
		tokenParams.Set("scope", "account+board")
		tokenParams.Set("grant_type", "authorization_code")
		tokenParams.Set("redirect_uri", redirectURL.String())
		tokenRequest, err := http.NewRequest("POST", "https://linuxfr.org/api/oauth/token", strings.NewReader(tokenParams.Encode()))
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot build linuxfr token request: %s", err), 500)
			return
		}
		tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		tokenRequest.Header.Set("User-Agent", "gb3")
		_, err = http.DefaultClient.Do(tokenRequest)
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot retrieve linuxfr token: %s", err), 500)
			return
		}
		http.Redirect(w, r, redirectURL.String(), 303)
	})
}
