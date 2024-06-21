package router

import (
	"fmt"
	"net"
)

type Response struct {
	conn       net.Conn
	statusCode int
	headers    map[string]string
	body       []byte
}

func (r *Response) WriteHeader(statusCode int) {
	r.statusCode = statusCode
}

func (r *Response) Header() map[string]string {
	return r.headers
}

func (r *Response) Write(body []byte) (int, error) {
	r.body = body
	return len(body), nil
}

func (r *Response) SendResponse() error {
	if r.statusCode == 0 {
		r.statusCode = 200
	}

	statusText := getStatusText(r.statusCode)
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.statusCode, statusText)

	for key, value := range r.headers {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	response += fmt.Sprintf("Content-Length: %d\r\n", len(r.body))
	response += "\r\n"

	_, err := r.conn.Write([]byte(response))
	if err != nil {
		return err
	}

	_, err = r.conn.Write(r.body)
	return err
}

func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 304:
		return "Not Modified"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown"
	}
}

func sendError(conn net.Conn, statusCode int, message string) {
	resp := &Response{
		conn:       conn,
		statusCode: statusCode,
		headers:    make(map[string]string),
	}
	resp.headers["Content-Type"] = "text/plain"
	resp.Write([]byte(message))
	resp.SendResponse()
}

func sendNotFound(resp *Response) {
	resp.WriteHeader(404)
	resp.headers["Content-Type"] = "text/plain"
	resp.Write([]byte("404 Not Found"))
	resp.SendResponse()
}
