package httpGateway

import (
	"github.com/smtdfc/photon/core"
	"github.com/smtdfc/photon/logger"
)

type Scope struct {
	logger      *logger.Logger
	prefix      string
	Gateway     *Gateway
	Middlewares []core.HttpHandler
}

func (s *Scope) Use(mw ...core.HttpHandler) {
	s.Middlewares = append(s.Middlewares, mw...)
}

func (s *Scope) SetLogger(logger *logger.Logger) {
	s.logger = logger
}

func (s *Scope) addRoute(method, pattern string, handlers ...core.HttpHandler) {
	fullPattern := s.prefix + pattern

	route := Route{
		Method:   method,
		Pattern:  fullPattern,
		Handlers: handlers,
		Scope:    s,
	}

	s.Gateway.Routes = append(s.Gateway.Routes, route)
}

func (s *Scope) Get(path string, handlers ...core.HttpHandler) {
	s.addRoute("GET", s.prefix+path, handlers...)
}

func (s *Scope) Post(path string, handlers ...core.HttpHandler) {
	s.addRoute("POST", s.prefix+path, handlers...)
}

func (s *Scope) Put(path string, handlers ...core.HttpHandler) {
	s.addRoute("PUT", s.prefix+path, handlers...)
}

func (s *Scope) Delete(path string, handlers ...core.HttpHandler) {
	s.addRoute("DELETE", s.prefix+path, handlers...)
}

func (s *Scope) Head(path string, handlers ...core.HttpHandler) {
	s.addRoute("HEAD", s.prefix+path, handlers...)
}

func (s *Scope) Patch(path string, handlers ...core.HttpHandler) {
	s.addRoute("PATCH", s.prefix+path, handlers...)
}

func (s *Scope) Options(path string, handlers ...core.HttpHandler) {
	s.addRoute("OPTIONS", s.prefix+path, handlers...)
}

func (s *Scope) Trace(path string, handlers ...core.HttpHandler) {
	s.addRoute("TRACE", s.prefix+path, handlers...)
}

func (s *Scope) Connect(path string, handlers ...core.HttpHandler) {
	s.addRoute("CONNECT", s.prefix+path, handlers...)
}
