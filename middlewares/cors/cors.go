package cors

import (
    "github.com/smtdfc/photon/core"
    "strings"
)

type Config struct {
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    AllowCredentials bool
}

func Middleware(cfg Config) func(ctx core.HttpContext) {
    return func(ctx core.HttpContext) {
        origin := ctx.GetHeader("Origin")

        allowOrigin := ""
        for _, o := range cfg.AllowedOrigins {
            if o == "*" || o == origin {
                allowOrigin = o
                break
            }
        }

        if allowOrigin != "" {
            ctx.SetHeader("Access-Control-Allow-Origin", allowOrigin)
            ctx.SetHeader("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
            ctx.SetHeader("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))

            if cfg.AllowCredentials {
                ctx.SetHeader("Access-Control-Allow-Credentials", "true")
            }
        }

        if ctx.Method() == "OPTIONS" {
            ctx.SetStatus(200)
            ctx.Abort()
            return
        }

        ctx.Next()
    }
}