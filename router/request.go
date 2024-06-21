package router

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type Request struct {
	Method      string
	Path        string
	Headers     map[string]string
	Body        []byte
	URLParams   map[string]string
	QueryParams map[string]string
}

func parseRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	parts := strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	method, path, _ := parts[0], parts[1], parts[2]

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	var body []byte
	if contentLength, ok := headers["Content-Length"]; ok {
		length := 0
		fmt.Sscanf(contentLength, "%d", &length)
		body = make([]byte, length)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return nil, err
		}
	}

	queryParams := make(map[string]string)
	if idx := strings.Index(path, "?"); idx != -1 {
		query := path[idx+1:]
		path = path[:idx]
		for _, param := range strings.Split(query, "&") {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) == 2 {
				queryParams[parts[0]] = parts[1]
			}
		}
	}

	return &Request{
		Method:      method,
		Path:        path,
		Headers:     headers,
		Body:        body,
		URLParams:   make(map[string]string),
		QueryParams: queryParams,
	}, nil
}
