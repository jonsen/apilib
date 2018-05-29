package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

// request header
type ReqHeader struct {
	Action string
	Time   string
	// other
}

// request
type Request struct {
	Header ReqHeader   // request header
	Body   interface{} // body
}

// response
type Response struct {
	Code    int         // custom status code
	Message string      // message
	Body    interface{} // body
}

// convert struct to string
func StructToString(data interface{}) string {
	byt, err := json.Marshal(data)
	if err != nil {
		return ""
	}

	return string(byt)
}

func (res *Response) String() string {
	str, err := json.Marshal(res)
	if err != nil {
		fmt.Println("error", err)
		return ""
	}

	fmt.Println(string(str))
	return string(str)
}

func (req *Request) String() string {
	str, err := json.Marshal(req)
	if err != nil {
		return ""
	}

	return string(str)
}

// read request data
func RequestReader(input io.ReadCloser, body interface{}) (req *Request, err error) {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("body is null")
	}

	req = new(Request)

	req.Body = body
	fmt.Println("RequestReader:", string(data))

	err = json.Unmarshal(data, req)

	return
}

// initial response data
// code: status code
// message: response messages
// body: response body
func ResponseWriter(code int, message string, body interface{}) (res *Response) {
	res = &Response{
		Code:    code,
		Message: message,
		Body:    body,
	}
	return
}
