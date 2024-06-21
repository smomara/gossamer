package template

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/smomara/gossamer/router"
)

var templates map[string]*template.Template

func InitTemplates(templateDir string) error {
	templates = make(map[string]*template.Template)

	files, err := os.ReadDir(templateDir)
	if err != nil {
		return fmt.Errorf("error reading template directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".html" {
			continue
		}

		tmpl, err := template.ParseFiles(filepath.Join(templateDir, file.Name()))
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", file.Name(), err)
		}

		templates[file.Name()] = tmpl
	}

	return nil
}

func RenderTemplate(w *router.Response, templateName string, data interface{}) error {
	tmpl, ok := templates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	w.Header()["Content-Type"] = "text/html"
	_, err = w.Write(buf.Bytes())
	return err
}
