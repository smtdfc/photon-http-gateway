package httpGateway

import (
	"github.com/smtdfc/photon/v2/core"
)

type Config struct {
	Port string
}

func New(config Config) *Gateway {
	return &Gateway{
		config:      config,
		Routes:      make([]Route, 0),
		Middlewares: make([]core.HttpHandler, 0),
	}
}
