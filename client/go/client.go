package client

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"strconv"
	//"fmt"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
	//"strconv"
)

const (
	version = "1.0"
	agent   = "APIClient Go client V1.0"
)

type APIClient struct {
	*http.Client
	LastTime string   // last time
	Proto    string   // http, https
	Host     string   // hostname
	Port     string   // port
	User     string   // user
	Passwd   string   // user's password
	AuthKey  string   // auth key
	SSLCa    string   // Ca certificate
	SSLCrt   string   // cert certificate
	SSLKey   string   // key certificate
	URL      *url.URL //
	Timeout  int      //
	Retry    int
	//
	SecureURI bool //
	XAPIKEY   bool //
	errorCode bool
}

type Options struct {
	Key       string // key
	Timeout   int    // timeout
	Retry     int    // retry
	SSLCa     string // Ca certificate
	SSLCrt    string // cert certificate
	SSLKey    string // key certificate
	SecureURI bool   //
	XAPIKEY   bool   //
	User      string // user
	Passwd    string // user's password
	AuthKey   string // auth key

}

// request header
type ReqHeader struct {
	Action interface{}
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

type RawResponse struct {
	StatusCode int
	ErrorCode  int
	Body       []byte
}

func (r RawResponse) Byte() []byte {
	return r.Body
}

func (r RawResponse) Response(v ...interface{}) (res *Response, err error) {
	res = new(Response)
	if len(v) > 0 && v[0] != nil {
		res.Body = v[0]
	}
	if r.Body == nil {
		return nil, errors.New("body is null")
	}
	d := json.NewDecoder(bytes.NewReader(r.Body))
	d.UseNumber()
	err = d.Decode(res)

	return
}

// encode md5
func Md5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

//
// example
//
// import (
//    "github.com/jonsen/apilib/client/go"
//    "fmt"
// )
// api, err := apilib.NewAPIClient("http://127.0.0.1:5188")
// OR
// opt := Options{}
// opt.Key = "test key"
// ...
// api, err := apilib.NewAPIClient("http://127.0.0.1:5188", opt)
// if err != nil {
//     fmt.Println(err)
// }
// // get test
// resp, err := api.Get("/acl/list")
// fmt.Println("GET", resp, err)
// [scheme:][//[userinfo@]host][/]path[?query][#fragment]
func NewAPIClient(rawurl string, options ...Options) (api *APIClient, err error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	api = new(APIClient)

	api.URL = u
	api.Host = u.Hostname()
	api.Port = u.Port()
	api.User = u.User.Username()
	api.Passwd, _ = u.User.Password()

	api.Retry = 3
	api.Timeout = 30

	if len(options) > 0 {
		if options[0].Retry > 0 {
			api.Retry = options[0].Retry
		}

		if options[0].Timeout > 0 {
			api.Timeout = options[0].Timeout
		}

		api.SecureURI = options[0].SecureURI
		api.XAPIKEY = options[0].XAPIKEY

		api.SSLCa = options[0].SSLCa
		api.SSLCrt = options[0].SSLCrt
		api.SSLKey = options[0].SSLKey
		if options[0].User != "" {
			api.User = options[0].User
		}
		if options[0].Passwd != "" {
			api.Passwd = options[0].Passwd
		}

		api.AuthKey = options[0].AuthKey
	}

	// set connect options
	dial := func(netw, addr string) (net.Conn, error) {
		deadline := time.Now().Add(30 * time.Second)
		c, err := net.DialTimeout(netw, addr, time.Second*time.Duration(api.Timeout))
		if err != nil {
			return nil, err
		}
		c.SetDeadline(deadline)
		return c, nil
	}

	// With SSL
	var tlsConfig *tls.Config
	if u.Scheme == "https" {

		// with CA
		if api.SSLCa != "" && api.SSLCrt != "" && api.SSLKey != "" {
			pool := x509.NewCertPool()

			caCrt, err := ioutil.ReadFile(api.SSLCa)
			if err != nil {
				return nil, err
			}
			pool.AppendCertsFromPEM(caCrt)

			cliCrt, err := tls.LoadX509KeyPair(api.SSLCrt, api.SSLKey)
			if err != nil {
				return nil, err
			}

			tlsConfig = &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			}
		} else {
			// no CA
			tlsConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}

	api.Client = &http.Client{
		Transport: &http.Transport{
			Dial:            dial,
			TLSClientConfig: tlsConfig,
		},
	}

	return
}

func (api *APIClient) SetRetry(retry int) {
	api.Retry = retry
}

func (api *APIClient) EnableSecureURI() {
	api.SecureURI = true
}

func (api *APIClient) EnableXAPIKEY() {
	api.XAPIKEY = true
}

func (api *APIClient) EnableErrorCode() {
	api.errorCode = true
}

// make request data
func makeRequest(action string, data ...interface{}) (b []byte, err error) {
	// header := ReqHeader{Action: action, Time: time.Now().Format("2006-01-02 15:04:05")}
	// req := &Request{Header: header, Body: data}
	header := ReqHeader{}
	if len(data) == 3 {
		header.Action = data[2]
		header.Time = time.Now().Format("2006-01-02 15:04:05")
	}

	req := &Request{Header: header, Body: data[0]}

	return json.Marshal(req)
}

// method: GET, POST, PUT, DELETE, PATCH,...
// url: query URL
// data: Is optional, When the GET method, is not available。
// data[0] Is to submit the data
// data[1] Is the data returned
// data[2] Is the header.Action
func (api *APIClient) method(method, url string, data ...interface{}) (raw RawResponse, err error) {
	var (
		req     *http.Request
		resp    *http.Response
		reqData []byte
		qURL    string
	)

	//
	if api.SecureURI || api.XAPIKEY {
		last := time.Now().Format("200601021504")
		if api.LastTime != last {
			api.LastTime = last
		}
	}
	if api.SecureURI {
		secureCode := Md5(api.User + ":" + api.Passwd + ":" + api.AuthKey + ":" + api.LastTime)
		qURL = "/" + secureCode + url
	} else {
		qURL = url
	}

	api.URL, err = api.URL.Parse(qURL)
	if err != nil {
		return
	}

	if len(data) < 1 {
		req, err = http.NewRequest(method, api.URL.String(), nil)
	} else {
		reqData, err = makeRequest(method, data...)
		if err == nil {
			req, err = http.NewRequest(method, api.URL.String(), bytes.NewReader(reqData))
		}
	}
	if err != nil {
		return
	}

	// set http head
	// create X-API-KEY
	if api.XAPIKEY {
		xapikey := Md5(qURL + ":" + string(reqData) + ":" + api.AuthKey + ":" + api.LastTime)
		req.Header.Set("X-API-KEY", xapikey)
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("User-Agent", agent)
	req.Header.Set("Connection", "close")

	// execution
	resp, err = api.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	raw.StatusCode = resp.StatusCode
	if api.errorCode {
		eCode := resp.Header.Get("Error-Code")
		raw.ErrorCode, _ = strconv.Atoi(eCode)
	}

	raw.Body, err = ioutil.ReadAll(resp.Body)

	return
}

func (api *APIClient) Method(method, url string, data ...interface{}) (raw RawResponse, err error) {
	for i := 0; i < api.Retry; i++ {
		raw, err = api.method(method, url, data...)
		if err == nil {
			break
		}
	}

	return
}

func (api *APIClient) Response(method, url string, data ...interface{}) (res *Response, err error) {
	raw, err := api.Method(method, url, data...)
	if err != nil {
		return
	}
	if len(data) < 2 {
		return raw.Response(nil)
	}

	return raw.Response(data[1])
}

// Get
func (api *APIClient) Get(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("GET", url, data...)
}

// Post
func (api *APIClient) Post(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("POST", url, data...)
}

// Put
func (api *APIClient) Put(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("PUT", url, data...)
}

// Patch
func (api *APIClient) Patch(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("PATCH", url, data...)
}

// Delete
func (api *APIClient) Delete(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("DELETE", url, data...)
}

// Head
func (api *APIClient) Head(url string) (raw RawResponse, err error) {
	return api.Method("HEAD", url)
}

// Connect
func (api *APIClient) Connect(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("CONNECT", url, data...)
}

// Trace
func (api *APIClient) Trace(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("TRACE", url, data...)
}

// Options
func (api *APIClient) Options(url string, data ...interface{}) (raw RawResponse, err error) {
	return api.Method("OPTIONS", url, data...)
}

//////
// GetResponse
func (api *APIClient) GetResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("GET", url, data...)
}

// PostResponse
func (api *APIClient) PostResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("POST", url, data...)
}

// PutResponse
func (api *APIClient) PutResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("PUT", url, data...)
}

// PatchResponse
func (api *APIClient) PatchResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("PATCH", url, data...)
}

// DeleteResponse
func (api *APIClient) DeleteResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("DELETE", url, data...)
}

// HeadResponse
func (api *APIClient) HeadResponse(url string) (res *Response, err error) {
	return api.Response("HEAD", url)
}

// ConnectResponse
func (api *APIClient) ConnectResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("CONNECT", url, data...)
}

// TraceResponse
func (api *APIClient) TraceResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("TRACE", url, data...)
}

// OptionsResponse
func (api *APIClient) OptionsResponse(url string, data ...interface{}) (res *Response, err error) {
	return api.Response("OPTIONS", url, data...)
}
