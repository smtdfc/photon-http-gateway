package httpGateway

import (
	"fmt"
	"github.com/smtdfc/photon/v2/core"
	"net/http"
)

type Gateway struct {
	Routes      []Route
	Middlewares []core.HttpHandler
	config      Config
}

func (g *Gateway) addRoute(method string, pattern string, handlers ...core.HttpHandler) {
	route := Route{
		Method:   method,
		Pattern:  pattern,
		Handlers: handlers,
		Scope:    nil,
	}
	g.Routes = append(g.Routes, route)
}

func (g *Gateway) findRoute(method string, rawPath string) (bool, Route, RouteParsed) {
	for _, route := range g.Routes {
		if route.Method != method {
			continue
		}

		ok, parsed := parseRoute(route.Pattern, rawPath)
		if ok {
			return true, route, parsed
		}
	}
	return false, Route{}, RouteParsed{}
}

func (g *Gateway) dispatch(route Route, ctx core.HttpContext, path string) {
	finalHandlers := append([]core.HttpHandler{}, g.Middlewares...)
	if route.Scope != nil {
		finalHandlers = append(finalHandlers, route.Scope.Middlewares...)
	}
	finalHandlers = append(finalHandlers, route.Handlers...)

	c := ctx.(*Context)
	scope := route.Scope
	if scope != nil && scope.logger != nil {
		scope.logger.Info(route.Method + " " + path)
	}

	c.handlers = finalHandlers
	c.index = -1

	c.Next()
}

func (g *Gateway) Use(mw ...core.HttpHandler) {
	g.Middlewares = append(g.Middlewares, mw...)
}

func (g *Gateway) CreateScope(module *core.Module, prefix string) core.HttpScope {
	return &Scope{
		Gateway: g,
		prefix:  prefix,
	}
}

func (g *Gateway) handler(w http.ResponseWriter, r *http.Request) {
	path := r.RequestURI
	method := r.Method
	passed, route, routeParsed := g.findRoute(method, path)
	if passed {
		ctx := NewContext(routeParsed, w, r)
		g.dispatch(route, ctx, path)
	}
}

func (g *Gateway) Start() error {
	var port string
	if g.config.Port != "" {
		port = g.config.Port
	} else {
		port = "3000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", g.handler)

	fmt.Println("Http server is running http://localhost:" + port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		fmt.Println("Error when starting Http server:", err)
		return err
	}
	return nil
}
