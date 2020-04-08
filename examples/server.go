package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	apilib "github.com/jonsen/apilib/server"
)

var (
	user = "user1"
	pass = "password2"
	key  = "testkey"
)

type hander struct {
}

func (h *hander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hander: ServeHTTP")
}

func fooHandler() http.Handler {
	return &hander{}
}

func main() {

	http.Handle("/foo", fooHandler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}

		defer r.Body.Close()

		// author basic
		if !apilib.AuthBasic(user, pass, r) {
			//apilib.AuthFailed(w)
			fmt.Println("AuthBasic AuthFailed")
			return
		}

		if !apilib.AuthXAPIKEY(key, r) {
			//apilib.AuthFailed(w)
			fmt.Println("AuthXAPIKEY AuthFailed")
			return
		}

		if !apilib.AuthSecureURI(user, pass, key, r) {
			//apilib.AuthFailed(w)
			fmt.Println("AuthXAPIKEY AuthFailed")
			return
		}

		fmt.Printf("request: %#v\n", r)
		fmt.Printf("query: %s, %s\n", r.Method, r.RequestURI)
		fmt.Printf("URL: %#v, %s\n", r.URL, r.URL.User.String())
		fmt.Printf("HEADER: %#v\n", r.Header)
		fmt.Printf("HEADER: X-API-KEY: %s\n", r.Header.Get("X-API-KEY"))

		reqData, err := apilib.RequestReader(r.Body, &data)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, apilib.ResponseWriter(8008, "read request body failed.", nil).String())
			return
		}
		fmt.Println("reqData", reqData)

		fmt.Fprintf(w, apilib.ResponseWriter(8000, "read request body success.", data).String())

	})

	go func() {
		err := apilib.Run(":5187", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err := apilib.Run(":5188", nil, "cert/server.crt", "cert/server.key", "cert/ca.crt")
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err := apilib.Run(":5189", nil, "cert/server.crt", "cert/server.key")
		if err != nil {
			fmt.Println(err)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)

	for {
		select {
		case s := <-ch:
			switch s {
			default:
			case syscall.SIGHUP:
				break
			case syscall.SIGINT:
				os.Exit(1)
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
			}
		}
	}
}
