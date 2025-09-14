package httpGateway

import (
	"github.com/smtdfc/photon/v2/core"
	"net/url"
	"strings"
)

type Route struct {
	Method   string
	Pattern  string
	Handlers []core.HttpHandler
	Scope    *Scope
}

type RouteParsed struct {
	Params map[string]string
	Query  map[string]string
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseRoute(pattern, rawPath string) (bool, RouteParsed) {
	u, err := url.Parse(rawPath)
	if err != nil {
		return false, RouteParsed{}
	}

	path := u.Path

	patternSplitted := strings.Split(pattern, "/")
	pathSplitted := strings.Split(path, "/")
	length := max(len(patternSplitted), len(pathSplitted))

	passed := true
	result := RouteParsed{
		Params: make(map[string]string),
		Query:  make(map[string]string),
	}

	for idx := 0; idx < length; idx++ {
		if idx >= len(patternSplitted) || idx >= len(pathSplitted) {
			passed = false
			break
		}

		if patternSplitted[idx] == "*" {
			break
		} else if strings.HasPrefix(patternSplitted[idx], ":") {
			runes := []rune(patternSplitted[idx])
			key := string(runes[1:])
			result.Params[key] = pathSplitted[idx]
		} else if patternSplitted[idx] != pathSplitted[idx] {
			passed = false
			break
		}
	}

	for k, v := range u.Query() {
		if len(v) > 0 {
			result.Query[k] = v[0]
		}
	}

	return passed, result
}
