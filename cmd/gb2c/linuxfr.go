package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type LinuxfrAccessData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Login        string `json:"login"`
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
		//Request token
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
		tokenResponse, err := gb2cHttpClient.Do(tokenRequest)
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot retrieve linuxfr token: %s", err), 500)
			return
		}

		//parse token
		var linuxfrAccessData LinuxfrAccessData
		defer tokenResponse.Body.Close()
		tokenDecoder := json.NewDecoder(tokenResponse.Body)
		err = tokenDecoder.Decode(&linuxfrAccessData)
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot decode linuxfr token: %s", err), 500)
			return
		}

		//retrieve login
		meResponse, err := gb2cHttpClient.Get(fmt.Sprintf("https://linuxfr.org/api/v1/me?bearer_token=%s", url.QueryEscape(linuxfrAccessData.AccessToken)))
		if nil != err {
			http.Error(w, fmt.Sprintf("linuxfr me request failed: %s", err), 500)
			return
		}
		defer meResponse.Body.Close()
		meDecoder := json.NewDecoder(meResponse.Body)
		err = meDecoder.Decode(&linuxfrAccessData)
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot decode linuxfr me response: %s", err), 500)
			return
		}

		//encode json access data (token and login) to base64 and send redirect
		body, err := json.Marshal(linuxfrAccessData)
		if nil != err {
			http.Error(w, fmt.Sprintf("Cannot encode linuxfr access data: %s", err), 500)
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
