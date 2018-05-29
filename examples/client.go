package main

import (
	"fmt"
	"log"

	apilib "github.com/jonsen/apilib/client/go"
)

func noneSSL() (api *apilib.APIClient, err error) {
	return apilib.NewAPIClient("http://user1:password1@localhost:5187")
}

func withCASSL() (api *apilib.APIClient, err error) {
	opt := apilib.Options{
		SSLCa:  "cert/ca.crt",
		SSLCrt: "cert/client.crt",
		SSLKey: "cert/client.key",
	}
	return apilib.NewAPIClient("https://user1:password2@localhost:5188", opt)
}

func noneCASSL() (api *apilib.APIClient, err error) {
	return apilib.NewAPIClient("https://user1:password3@localhost:5189")
}

func action(api *apilib.APIClient) {
	api.EnableSecureURI()

	api.EnableXAPIKEY()

	resp, err := api.Get("/test/get?mod=test&act=ok")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(api)

	fmt.Println(resp)

	data := map[string]string{"mod": "post", "act": "true"}
	resp, err = api.Post("/test/post", data)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(api)
	fmt.Println(resp)

	resp, err = api.Delete("/test/delete", data)
	if err != nil {
		log.Println(err)
		return
	}
	resp, err = api.Put("/test/put", data)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err = api.Head("/test/head")
	if err != nil {
		log.Println(err)
		return
	}
}

func main() {
	api1, err := noneSSL()
	if err != nil {
		log.Fatalf("noneSSL:%s", err)
	}
	action(api1)
	api2, err := noneCASSL()
	if err != nil {
		log.Fatal("noneCASSL:%s", err)
	}
	action(api2)
	api3, err := withCASSL()
	if err != nil {
		log.Fatal("withCASSL:%s", err)
	}
	action(api3)
}
