package httpGateway

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/smtdfc/photon/core"
)

type Context struct {
	req         *http.Request
	writer      http.ResponseWriter
	routeParsed RouteParsed
	aborted     bool
	index       int
	data        map[string]any
	mu          sync.RWMutex
	handlers    []core.HttpHandler
}

func NewContext(routeParsed RouteParsed, writer http.ResponseWriter, req *http.Request) core.HttpContext {
	return &Context{
		req:         req,
		writer:      writer,
		routeParsed: routeParsed,
		data:        make(map[string]any),
	}
}

func (c *Context) Method() string   { return c.req.Method }
func (c *Context) Path() string     { return c.req.URL.Path }
func (c *Context) Protocol() string { return c.req.Proto }

func (c *Context) Param(key string) string {
	if c.routeParsed.Params[key] != "" {
		return c.routeParsed.Params[key]
	}
	return ""
}

func (c *Context) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

func (c *Context) QueryDefault(key, def string) string {
	if v := c.req.URL.Query().Get(key); v != "" {
		return v
	}
	return def
}

func (c *Context) Header(key string) string {
	return c.req.Header.Get(key)
}

func (c *Context) Cookie(name string) string {
	ck, err := c.req.Cookie(name)
	if err != nil {
		return ""
	}
	return ck.Value
}

func (c *Context) Body() []byte {
	b, _ := io.ReadAll(c.req.Body)
	return b
}

func (c *Context) FormValue(key string) string {
	return c.req.FormValue(key)
}

func (c *Context) FormFile(name string) ([]byte, error) {
	file, _, err := c.req.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

func (c *Context) Status(code int) core.HttpContext {
	c.writer.WriteHeader(code)
	return c
}

func (c *Context) SetHeader(key, value string) core.HttpContext {
	c.writer.Header().Set(key, value)
	return c
}

func (c *Context) SetCookie(name, value string, options ...any) core.HttpContext {
	cookie := &http.Cookie{Name: name, Value: value}
	http.SetCookie(c.writer, cookie)
	return c
}

func (c *Context) Text(code int, data string) core.HttpContext {
	c.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.writer.WriteHeader(code)
	_, _ = c.writer.Write([]byte(data))
	return c
}

func (c *Context) JSON(code int, data any) core.HttpContext {
	c.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.writer.WriteHeader(code)
	_ = json.NewEncoder(c.writer).Encode(data)
	return c
}

func (c *Context) HTML(code int, html string) core.HttpContext {
	c.writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.writer.WriteHeader(code)
	_, _ = c.writer.Write([]byte(html))
	return c
}

func (c *Context) Blob(code int, contentType string, data []byte) core.HttpContext {
	c.writer.Header().Set("Content-Type", contentType)
	c.writer.WriteHeader(code)
	_, _ = c.writer.Write(data)
	return c
}

func (c *Context) File(code int, filepath string) core.HttpContext {
	f, err := os.Open(filepath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return c
	}
	defer f.Close()

	buff := make([]byte, 512)
	n, _ := f.Read(buff)
	contentType := http.DetectContentType(buff[:n])

	c.writer.Header().Set("Content-Type", contentType)
	c.writer.WriteHeader(code)
	_, _ = c.writer.Write(buff[:n])
	_, _ = io.Copy(c.writer, f)
	return c
}

func (c *Context) Next() core.HttpContext {
	c.index++
	for c.index < len(c.handlers) {
		if c.aborted {
			break
		}
		c.handlers[c.index](c)
		c.index++
	}
	return c
}

func (c *Context) Abort() core.HttpContext {
	c.aborted = true
	return c
}

func (c *Context) IsAborted() bool {
	return c.aborted
}

func (c *Context) Set(key string, value any) core.HttpContext {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return c
}

func (c *Context) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key]
}

func (c *Context) MustGet(key string) any {
	if v := c.Get(key); v != nil {
		return v
	}
	panic("key not found: " + key)
}
