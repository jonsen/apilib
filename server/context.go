package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

// Request read client's json data
func (s *Server) Request(c *gin.Context, body interface{}) (req *Request, err error) {
	defer c.Request.Body.Close()

	req, err = RequestReader(c.Request.Body, body)
	return
}

// Response write json to client
func (s *Server) Response(c *gin.Context, code int, body interface{}, message string, v ...interface{}) {
	c.JSON(200, ResponseWriter(code, fmt.Sprintf(message, v...), body))
}
