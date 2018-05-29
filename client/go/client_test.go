package client

import (
	"fmt"
	"testing"
)

func TestAPILib(t *testing.T) {
	/*
		opt := Options{
			SSLCa:  "cert/ca.crt",
			SSLCrt: "cert/client.crt",
			SSLKey: "cert/client.key",
		}
	*/
	opt := Options{}
	api, err := NewAPIClient("http://user1:password2@127.0.0.1", opt)
	if err != nil {
		t.Fatal(err)
	}

	//api.EnableSecureURI()

	api.EnableXAPIKEY()

	raw, err := api.Get("/test/echojson.php?mod=test&act=ok")
	if err != nil {
		//t.Fatal(err)
	}
	fmt.Println(api, (raw))

	data := map[string]string{"mod": "post", "act": "true"}
	raw, err = api.Post("/test/post", data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(api)
	fmt.Println(raw)

	raw, err = api.Delete("/test/delete", data)
	if err != nil {
		t.Fatal(err)
	}
	raw, err = api.Put("/test/put", data)
	if err != nil {
		t.Fatal(err)
	}

	raw, err = api.Head("/test/head")
	if err != nil {
		t.Fatal(err)
	}
}
