package slick

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/julienschmidt/httprouter"
)

type Plug func(Handler) Handler

type Handler func(c *Context) error

type ErrorHandler func(error, *Context) error

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	ctx      context.Context
	params   httprouter.Params
}

func newContext(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	return &Context{
		Response: w,
		Request:  r,
		ctx:      context.Background(),
		params:   params,
	}
}

func (c *Context) Param(name string) string {
	return c.params.ByName(name)
}

func (c *Context) Query(name string) string {
	return c.Request.URL.Query().Get(name)
}

func (c *Context) FormValue(name string) string {
	return c.Request.FormValue(name)
}

func (c *Context) Render(component templ.Component) error {
	return component.Render(c.ctx, c.Response)
}

func (c *Context) Redirect(url string, code int) error {
	if code < http.StatusMultipleChoices || code > http.StatusTemporaryRedirect {
		return errors.New("invalid redirect code")
	}
	http.Redirect(c.Response, c.Request, url, code)
	return nil
}

func (c *Context) JSON(status int, v any) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	return json.NewDecoder(c.Request.Body).Decode(&v)
}

func (c *Context) Text(status int, t string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(t))
	return err
}

func (c *Context) Set(key string, value any) {
	c.ctx = context.WithValue(c.ctx, key, value)
}

func (c *Context) Get(key string) any {
	return c.ctx.Value(key)
}

type Slick struct {
	ErrorHandler ErrorHandler
	router       *httprouter.Router
	plugs        []Plug
}

func New() *Slick {
	return &Slick{
		router:       httprouter.New(),
		ErrorHandler: defaultErrorHandler,
	}
}

func (s *Slick) HandleMethodNotAllowed(h http.Handler) {
	s.router.MethodNotAllowed = h
}

func (s *Slick) Plug(plugs ...Plug) {
	s.plugs = append(s.plugs, plugs...)
}

func (s *Slick) Start(port string) error {
	fmt.Printf("slick app running http://localhost:%s\n", port)
	return http.ListenAndServe(port, s.router)
}

func (s *Slick) add(method, path string, h Handler, plugs ...Plug) {
	s.router.Handle(method, path, s.makeHTTPRouterHandle(h, plugs...))
}

func (s *Slick) Get(path string, h Handler, plugs ...Plug) {
	s.add("GET", path, h, plugs...)
}

func (s *Slick) Post(path string, h Handler, plugs ...Plug) {
	s.add("POST", path, h, plugs...)
}

func (s *Slick) Put(path string, h Handler, plugs ...Plug) {
	s.add("PUT", path, h, plugs...)
}

func (s *Slick) Delete(path string, h Handler, plugs ...Plug) {
	s.add("DELETE", path, h, plugs...)
}

func (s *Slick) Head(path string, h Handler, plugs ...Plug) {
	s.add("HEAD", path, h, plugs...)
}

func (s *Slick) Options(path string, h Handler, plugs ...Plug) {
	s.add("OPTIONS", path, h, plugs...)
}

func (s *Slick) makeHTTPRouterHandle(h Handler, plugs ...Plug) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := newContext(w, r, params)
		for i := len(plugs) - 1; i >= 0; i-- {
			h = plugs[i](h)
		}
		for i := len(s.plugs) - 1; i >= 0; i-- {
			h = s.plugs[i](h)
		}
		if err := h(ctx); err != nil {
			// todo: handle the error from the error handler huh?
			s.ErrorHandler(err, ctx)
		}
	}
}

func defaultErrorHandler(err error, c *Context) error {
	slog.Error("error", "err", err)
	return nil
}
