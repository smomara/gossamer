package router

import (
	"net"
	"strings"
)

type Handler func(w *Response, r *Request)

type Route struct {
	Method  string
	Path    string
	Handler Handler
}

type Router struct {
	routes []Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) AddRoute(method, path string, handler Handler) {
	r.routes = append(r.routes, Route{Method: method, Path: path, Handler: handler})
}

func (r *Router) ServeHTTP(conn net.Conn) error {
	defer conn.Close()

	req, err := parseRequest(conn)
	if err != nil {
		sendError(conn, 400, "Bad Request")
		return err
	}

	resp := &Response{conn: conn, headers: make(map[string]string)}

	for _, route := range r.routes {
		if route.Method == req.Method && matchPath(route.Path, req.Path) {
			req.URLParams = extractURLParams(route.Path, req.Path)
			route.Handler(resp, req)
			return nil
		}
	}
	sendNotFound(resp)
	return nil
}

func matchPath(routePath, requestPath string) bool {
	routeParts := strings.Split(strings.Trim(routePath, "/"), "/")
	requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")

	if len(routeParts) > 0 && routeParts[len(routeParts)-1] == "*" {
		return strings.HasPrefix(requestPath, strings.TrimSuffix(routePath, "/*"))
	}

	if len(routeParts) != len(requestParts) {
		return false
	}

	for i, routePart := range routeParts {
		if routePart == "*" || strings.HasPrefix(routePart, ":") {
			continue
		}
		if routePart != requestParts[i] {
			return false
		}
	}

	return true
}

func extractURLParams(routePath, requestPath string) map[string]string {
	params := make(map[string]string)

	routeParts := strings.Split(strings.Trim(routePath, "/"), "/")
	requestParts := strings.Split(strings.Trim(requestPath, "/"), "/")

	for i, routePart := range routeParts {
		if strings.HasPrefix(routePart, ":") {
			paramName := strings.TrimPrefix(routePart, ":")
			params[paramName] = requestParts[i]
		}
	}

	if routeParts[len(routeParts)-1] == "*" {
		params["*"] = strings.Join(requestParts[len(routeParts)-1:], "/")
	}

	return params
}
