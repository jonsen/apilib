package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Run(addr string, m http.Handler, ssl ...string) error {
	l := len(ssl)
	if l == 3 {
		if ssl[2] != "" {
			pool := x509.NewCertPool()

			caCrt, err := ioutil.ReadFile(ssl[2])
			if err != nil {
				panic("ReadFile CA err:" + err.Error())
			}
			pool.AppendCertsFromPEM(caCrt)

			s := &http.Server{
				Addr:    addr,
				Handler: m,
				TLSConfig: &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.RequireAndVerifyClientCert,
				},
			}

			return s.ListenAndServeTLS(ssl[0], ssl[1])
		} else {
			return http.ListenAndServeTLS(addr, ssl[0], ssl[1], m)
		}
	}
	if l == 2 {
		return http.ListenAndServeTLS(addr, ssl[0], ssl[1], m)
	}

	return http.ListenAndServe(addr, m)
}

func GetValues(requestURI string) (values map[string]string, err error) {
	u, err := url.Parse(requestURI)
	if err != nil {
		return
	}

	uValues := u.Query()
	for k, v := range uValues {
		if len(v) > 0 {
			values[k] = v[0]
		} else {
			values[k] = ""
		}
	}

	return
}
