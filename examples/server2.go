package main

import (
	"fmt"
	//"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonsen/apilib/server"
)

var (
	user = "user1"
	pass = "password2"
	key  = "testkey"
)

func main() {

	svr := server.NewServer("apiServer", "1.0")
	//svr.AuthBasic("user1", "pass1")

	svr.Get("/foo", func(c *server.Context) {
		c.Text(200, "foo func")

	})

	svr.Get("/json", func(c *server.Context) {
		c.Response(200, map[string]interface{}{"user": "my name is xxx"}, "ok")
	})

	svr.Post("/echo", func(c *server.Context) {
		var body map[string]interface{}
		req, err := c.Request(&body)
		if err != nil {
			fmt.Println(err)
			return
		}

		c.Response(200, req, "ok")
	})

	go func() {
		err := svr.Run(":5197")
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err := svr.Run(":5198", "cert/server.crt", "cert/server.key", "cert/ca.crt")
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		err := svr.Run(":5199", "cert/server.crt", "cert/server.key")
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
