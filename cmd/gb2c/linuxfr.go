package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func remoteHost(r *http.Request) string {
	host := r.Header.Get("X-Forwarded-Host")
	if len(host) == 0 {
		host = r.Host
	}
	return host
}

//RegisterLinuxfrAPI setup http linuxfr API endpoints
func RegisterLinuxfrAPI() {
	clientID := os.Getenv("GB2C_LINUXFR_CLIENT_ID")
	clientSecret := os.Getenv("GB2C_LINUXFR_CLIENT_SECRET")
	if len(clientID) == 0 || len(clientSecret) == 0 {
		log.Println("Do not forget to GB2C_LINUXFR_CLIENT_ID and GB2C_LINUXFR_CLIENT_SECRET environment variable if you want to post to linuxfr tribune")
		return
	}
	http.HandleFunc("/gb2c/linuxfr/authorize", func(w http.ResponseWriter, r *http.Request) {
		linuxfrRedirectURL := url.URL{
			Scheme: "https",
			Host:   remoteHost(r),
			Path:   "/gb2c/linuxfr/connected",
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
		http.Redirect(w, r, linuxfrAuthorizeURL.String(), http.StatusSeeOther)
	})
	http.HandleFunc("/gb2c/linuxfr/connected", func(w http.ResponseWriter, r *http.Request) {

		linuxfrRedirectURL := url.URL{
			Scheme: "https",
			Host:   remoteHost(r),
			Path:   "/gb2c/linuxfr/connected",
		}

		tokenParams := url.Values{}
		tokenParams.Set("client_id", clientID)
		tokenParams.Set("client_secret", clientSecret)
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
			Host:     remoteHost(r),
			Path:     "/linuxfr.html",
			Fragment: encodedToken,
		}
		http.Redirect(w, r, redirectURL.String(), http.StatusSeeOther)
	})
}
