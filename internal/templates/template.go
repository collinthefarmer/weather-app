package templates

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type TemplateEngine struct {
	root *template.Template
}

func Init(fs embed.FS, rootTemplatePath string, commonTemplatePaths ...string) (*TemplateEngine, error) {
	root, err := template.ParseFS(fs, rootTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing root template: %w", err)
	}

	templateFunctions := template.FuncMap{}

	return &TemplateEngine{root: template.Must(
		root.
			Funcs(templateFunctions).
			ParseFS(fs, commonTemplatePaths...),
	)}, nil
}

func (te *TemplateEngine) Render(w http.ResponseWriter, path string, data any) error {
	tmpl := template.Must(te.root.Clone())
	tmpl, err := tmpl.ParseFiles(path)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}
