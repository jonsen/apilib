package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

type (
	Server struct {
		*martini.Martini
		Router
	}

	// sessions
	Session   sessions.Session
	Render    render.Render
	RenderOpt render.Options
)

func NewServer(app, version string) *Server {
	r := NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(NewContext(CtxOptions{AppName: app, Version: version}))
	m.MapTo(r, (*Routes)(nil))
	m.Action(r.Handle)

	return &Server{m, r}
}

func (s *Server) Run(addr string, ssl ...string) error {
	return Run(addr, s, ssl...)
}

func (s *Server) AuthBasic(user, pass string) {
	s.Use(func(c *Context) {
		if !AuthBasic(user, pass, c.Req) {
			c.AuthFailed()
			return
		}

		c.Next()
	})
}

func (s *Server) AuthXAPIKEY(key string) {
	s.Use(func(c *Context) {
		if !AuthXAPIKEY(key, c.Req) {
			c.AuthFailed()
			return
		}

		c.Next()
	})
}

func (s *Server) AuthSecureURI(user, pass, key string) {
	s.EnableSecure()
	s.Use(func(c *Context) {
		if !AuthSecureURI(user, pass, key, c.Req) {
			c.AuthFailed()
			return
		}

		c.Next()
	})
}

func (s *Server) AuthClient(allows []string) {
	s.Use(func(c *Context) {
		if !AuthClient(allows, c.Req) {
			c.AuthFailed()
			return
		}

		c.Next()
	})
}

///////

func (s *Server) UseSession(key string) {
	store := sessions.NewCookieStore([]byte(key))
	s.Use(sessions.Sessions("api_session", store))
}

////
func (s *Server) UseRender() {
	s.Use(render.Renderer())
}

func (s *Server) NotFount() {
	s.NotFound(func(c *Context) {
		c.NotFound()
	})
}

// Static server
func (s *Server) Static(path, uri string) {
	s.Use(martini.Static(path, martini.StaticOptions{Prefix: uri}))
}
