package static

import (
	"log"
	"mime"
	"os"
	"path/filepath"

	"github.com/smomara/gossamer/router"
)

func ServeStaticFiles(r *router.Router, urlPrefix, dirPath string) {
	r.AddRoute("GET", urlPrefix+"/*", func(w *router.Response, r *router.Request) {
		requestedPath := r.URLParams["*"]
		filePath := filepath.Join(dirPath, requestedPath)

		log.Printf("Requested static file: %s\n", filePath)

		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			log.Printf("File not found: %s\n", filePath)
			w.WriteHeader(404)
			w.Write([]byte("404 Not Found"))
			return
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file: %s, Error: %v\n", filePath, err)
			w.WriteHeader(500)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header()["Content-Type"] = contentType
		w.Header()["X-Content-Type-Options"] = "nosniff"
		w.WriteHeader(200)

		log.Printf("Serving file: %s with content type: %s\n", filePath, contentType)

		w.Write(content)
		w.SendResponse()
	})
}
