package main

import (
	"fmt"
	//"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jonsen/apilib/server"
)

var (
	user = "user1"
	pass = "password2"
	key  = "testkey"
)

func main() {

	svr := server.NewServer("apiServer", "1.0", "debug")
	//svr.AuthBasic("user1", "pass1")

	svr.GET("/foo", func(c *gin.Context) {
		c.String(200, "foo func")

	})

	svr.GET("/json", func(c *gin.Context) {
		svr.Response(c, 200, map[string]interface{}{"user": "my name is xxx"}, "ok")
	})

	svr.POST("/echo", func(c *gin.Context) {
		var body map[string]interface{}
		req, err := svr.Request(c, &body)
		if err != nil {
			fmt.Println(err)
			return
		}

		svr.Response(c, 200, req, "ok")
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
