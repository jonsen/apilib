package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"net/url"
)

func RunTLS(addr, ca, crt, key string, m http.Handler) error {

	if ca == "" {
		return http.ListenAndServeTLS(addr, crt, key, m)
	}

	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile(ca)
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

	return s.ListenAndServeTLS(crt, key)
}

func Run(addr string, m http.Handler) error {

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
