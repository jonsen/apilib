package server

import (
	"github.com/gin-gonic/gin"
)

type (
	Server struct {
		*gin.Engine
		Application string
		Version     string
	}
)

func NewServer(app, version, mode string) *Server {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())
	gin.SetMode(mode)

	return &Server{e, app, version}
}

func (s *Server) Run(addr string, ssl ...string) error {
	return Run(addr, s, ssl...)
}

func (s *Server) AuthBasic(user, pass string) {
	s.Use(func(c *gin.Context) {
		if !AuthBasic(user, pass, c.Request) {
			s.AuthFailed(c)
			return
		}

		c.Next()
	})
}

func (s *Server) AuthXAPIKEY(key string) {
	s.Use(func(c *gin.Context) {
		if !AuthXAPIKEY(key, c.Request) {
			s.AuthFailed(c)
			return
		}

		c.Next()
	})
}

func (s *Server) AuthSecureURI(user, pass, key string) {
	//s.EnableSecure()
	s.Use(func(c *gin.Context) {
		if !AuthSecureURI(user, pass, key, c.Request) {
			s.AuthFailed(c)
			return
		}

		c.Next()
	})
}

func (s *Server) AuthClient(allows []string) {
	s.Use(func(c *gin.Context) {
		if !AuthClient(allows, c.Request) {
			s.AuthFailed(c)
			return
		}

		c.Next()
	})
}

// AuthFailed write auth failed message to client
func (s *Server) AuthFailed(c *gin.Context) {
	c.JSON(403, ResponseWriter(403, "authorization failed", nil))
}

// NotFound write page not found message to client
func (s *Server) NotFound(c *gin.Context) {
	c.JSON(404, ResponseWriter(404, "page not found", nil))
}

///////

// func (s *Server) UseSession(key string) {
// 	store := sessions.NewCookieStore([]byte(key))
// 	s.Use(sessions.Sessions("api_session", store))
// }

////
// func (s *Server) UseRender() {
// 	s.Use(render.Renderer())
// }

// func (s *Server) NotFount() {
// 	s.NotFound(func(c *Context) {
// 		c.NotFound()
// 	})
// }

// Static server
// func (s *Server) Static(path, uri string) {
// 	s.Use(martini.Static(path, martini.StaticOptions{Prefix: uri}))
// }
