package server

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// encode md5
func Md5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

func AuthFailed(w http.ResponseWriter) {
	http.Error(w, "authorization failed", http.StatusUnauthorized)
}

func AuthBasic(user, pass string, r *http.Request) bool {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		return false
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 || pair[0] != user || pair[1] != pass {
		return false
	}

	return true
}

func AuthXAPIKEY(key string, r *http.Request) bool {
	last := time.Now().Format("200601021504")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false
	}

	lkey := Md5(r.RequestURI + ":" + string(reqBody) + ":" + key + ":" + last)
	rkey := r.Header.Get("X-API-KEY")

	// check X-API-KEY
	if lkey != rkey {
		return false
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

	return true
}

func AuthSecureURI(user, pass, key string, r *http.Request) bool {
	// check URI's code
	last := time.Now().Format("200601021504")
	paths := strings.SplitN(r.RequestURI, "/", 3)
	if len(paths) < 2 || paths[1] != Md5(user+":"+pass+":"+key+":"+last) {
		return false
	}

	return true
}

func AuthClient(allows []string, r *http.Request) (auth bool) {
	auth = false

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return
	}
	rAddr := net.ParseIP(host)

	//hosts := strings.Split(allows, ";")
	for _, v := range allows {
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			ipHost := net.ParseIP(v)
			if ipHost != nil {
				if ipHost.Equal(rAddr) {
					auth = true
					break
				}
			} else {
			}
		} else {
			if ipNet.Contains(rAddr) {
				auth = true
				break
			}
		}
	}
	return
}
