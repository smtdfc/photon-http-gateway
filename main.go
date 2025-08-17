package httpGateway

import (
	"fmt"
	"net/http"
	"github.com/smtdfc/photon/core"
)

type Config struct {
	Port string
}

type Scope struct{
	Gateway *Gateway
	Middlewares []core.HttpHandler
}

func (s *Scope) Use(mw ...core.HttpHandler){
	
}

func (s *Scope) Get(path string, handlers ...core.HttpHandler){
	
}

func (s *Scope) Post(path string, handlers ...core.HttpHandler){
	
}

func (s *Scope) Put(path string, handlers ...core.HttpHandler){
	
}

func (s *Scope) Delete(path string, handlers ...core.HttpHandler){
	
}

func (s *Scope) Head(path string, handlers ...core.HttpHandler){
	
}

type Gateway struct {
	Middlewares []core.HttpHandler
	config Config
}

func (g *Gateway) Use(mw ...core.HttpHandler){
	
}

func (g *Gateway) CreateScope(module *core.Module, prefix string) core.HttpScope{
	return &Scope{
		Gateway:g,
	}
}

func (g *Gateway) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func New(config Config) *Gateway {
	return &Gateway{
		config:config,
	}
}

func (g *Gateway) Start() error {
	var port string 
	if g.config.Port != "" {
		port = g.config.Port
	}else{
		port = "3000"
	}

	http.HandleFunc("/", g.handler)
	fmt.Println("Http server is running http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error when starring Http server:", err)
		return err
	}
	return nil
}
