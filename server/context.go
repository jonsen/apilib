package server

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	//"github.com/codegangsta/inject"
	"github.com/go-martini/martini"
	//"github.com/martini-contrib/render"
)

const (
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	ContentBinary  = "application/octet-stream"
	ContentText    = "text/plain; charset=UTF-8"
	ContentJSON    = "application/json; charset=UTF-8"
	ContentHTML    = "text/html; charset=UTF-8"
	ContentXHTML   = "application/xhtml+xml; charset=UTF-8"
	ContentXML     = "text/xml; charset=UTF-8"
	defaultCharset = "UTF-8"
)

//type Params martini.Params

type Context struct {
	martini.Context
	http.ResponseWriter
	Req  *http.Request
	data *sync.Map
	opt  CtxOptions
}

/*
type context struct {
	inject.Injector
	handlers []Handler
	action   Handler
	rw       ResponseWriter
	index    int
}
*/
type CtxOptions struct {
	AppName string
	Version string
	// Appends the given charset to the Content-Type header. Default is "UTF-8".
	Charset string
	// Outputs human readable JSON
	IndentJSON bool
	// Outputs human readable XML
	IndentXML bool
	// Prefixes the JSON output with the given bytes.
	PrefixJSON []byte
	// Prefixes the XML output with the given bytes.
	PrefixXML []byte
	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
	HTMLContentType string
}

// NewContext new context
func NewContext(options ...CtxOptions) martini.Handler {
	data := new(sync.Map)

	var opt CtxOptions
	if len(options) > 0 {
		opt = options[0]
	}
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		c.Map(&Context{c, res, req, data, opt})
	}
}

// Get get value by key
func (p Params) Get(key string) string {
	return p[key]
}

// Params get all params from router
// exmaple
//    s.Get("/user/:id", getUser)
//
//    func getUser(c *Content) {
//       params := c.Params()
//       id := params["id"]
//     // or
//       id := params.Get("id")
//     // or
//       id := c.Params().Get("id")
//    }
//
func (c *Context) Params() (params Params) {
	var pm martini.Params
	value := c.Context.Get(reflect.TypeOf(pm))
	if !value.IsValid() {
		return
	}
	if value.Type().Kind() != reflect.Map {
		return
	}
	new := value.Convert(reflect.TypeOf(params))
	if !new.IsValid() {
		return
	}

	params = new.Interface().(Params)

	return
}

// JSON render JSON
func (c *Context) JSON(status int, v interface{}) {
	var result []byte
	var err error
	if c.opt.IndentJSON {
		result, err = json.MarshalIndent(v, "", "  ")
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		http.Error(c, err.Error(), 500)
		return
	}

	// json rendered fine, write out the result
	c.Header().Set(ContentType, ContentJSON)
	c.WriteHeader(status)
	if len(c.opt.PrefixJSON) > 0 {
		c.Write(c.opt.PrefixJSON)
	}
	c.Write(result)
}

// Data render raw data
func (c *Context) Data(status int, v []byte) {
	if c.Header().Get(ContentType) == "" {
		c.Header().Set(ContentType, ContentBinary)
	}
	c.WriteHeader(status)
	c.Write(v)
}

// Text render Text
func (c *Context) Text(status int, v string) {
	if c.Header().Get(ContentType) == "" {
		c.Header().Set(ContentType, ContentText)
	}
	c.WriteHeader(status)
	c.Write([]byte(v))
}

// XML render XML
func (c *Context) XML(status int, v interface{}) {
	var result []byte
	var err error
	if c.opt.IndentXML {
		result, err = xml.MarshalIndent(v, "", "  ")
	} else {
		result, err = xml.Marshal(v)
	}
	if err != nil {
		http.Error(c, err.Error(), 500)
		return
	}

	// XML rendered fine, write out the result
	c.Header().Set(ContentType, ContentXML)
	c.WriteHeader(status)
	if len(c.opt.PrefixXML) > 0 {
		c.Write(c.opt.PrefixXML)
	}
	c.Write(result)
}

// GetData get data
func (c *Context) GetData(key string) (interface{}, bool) {
	return c.data.Load(key)
}

// SetData set data
func (c *Context) SetData(key string, value interface{}) {
	if c.data == nil {
		c.data = new(sync.Map)
	}
	c.data.Store(key, value)
}

// DelData delete data
func (c *Context) DelData(key string) {
	c.data.Delete(key)
}

// Request read client's data
func (c *Context) Request(body interface{}) (req *Request, err error) {
	defer c.Req.Body.Close()

	req, err = RequestReader(c.Req.Body, body)
	return
}

// Response write data to client
func (c *Context) Response(code int, body interface{}, message string, v ...interface{}) {
	c.JSON(200, ResponseWriter(code, fmt.Sprintf(message, v...), body))
}

// AuthFailed write auth failed message to client
func (c *Context) AuthFailed() {
	c.JSON(403, ResponseWriter(403, "authorization failed", nil))
}

// NotFound write page not found message to client
func (c *Context) NotFound() {
	c.JSON(404, ResponseWriter(404, "page not found", nil))
}
