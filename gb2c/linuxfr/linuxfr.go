package linuxfr

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
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
	if len(clientID) == 0 {
		return
	}
	http.HandleFunc("/api/linuxfr/authorize", func(w http.ResponseWriter, r *http.Request) {
		linuxfrRedirectURL := url.URL{
			Scheme: "https",
			Host:   r.Host,
			Path:   "/api/linuxfr/connected",
		}

		linuxfrAuthorizeQuery := make(url.Values)
		linuxfrAuthorizeQuery.Set("client_id", clientID)
		linuxfrAuthorizeQuery.Set("response_type", "code")
		linuxfrAuthorizeQuery.Set("scope", "account board")
		linuxfrAuthorizeQuery.Set("redirect_uri", linuxfrRedirectURL.String())
		linuxfrAuthorizeURL := url.URL{
			Scheme:   "https",
			Host:     "linuxfr.org",
			Path:     "/api/oauth/authorize",
			RawQuery: linuxfrAuthorizeQuery.Encode(),
		}
		http.Redirect(w, r, linuxfrAuthorizeURL.String(), 303)
	})
	http.HandleFunc("/api/linuxfr/connected", func(w http.ResponseWriter, r *http.Request) {
		linuxfrRedirectURL := url.URL{
			Scheme: "https",
			Host:   r.Host,
			Path:   "/api/linuxfr/connected",
		}

		tokenParams := url.Values{}
		tokenParams.Set("client_id", clientID)
		tokenParams.Set("code", r.URL.Query().Get("code"))
		tokenParams.Set("grant_type", "authorization_code")
		tokenParams.Set("redirect_uri", linuxfrRedirectURL.String())
		tokenRequest, err := http.NewRequest("POST", "https://linuxfr.org/api/oauth/token", strings.NewReader(tokenParams.Encode()))
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot build linuxfr token request: %s", err), 500)
			return
		}
		tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		tokenRequest.Header.Set("User-Agent", "gb3")
		tokenResponse, err := http.DefaultClient.Do(tokenRequest)
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot retrieve linuxfr token: %s", err), 500)
			return
		}
		defer tokenResponse.Body.Close()
		body, err := ioutil.ReadAll(tokenResponse.Body)
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot read linuxfr token body: %s", err), 500)
			return
		}
		encodedToken := base64.StdEncoding.EncodeToString(body)
		redirectURL := url.URL{
			Scheme:   "https",
			Host:     r.Host,
			Path:     "/linuxfr.html",
			Fragment: encodedToken,
		}
		http.Redirect(w, r, redirectURL.String(), 303)
	})
}
