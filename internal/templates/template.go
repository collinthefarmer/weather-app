package templates

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type TemplateEngine struct {
	root      *template.Template
	constants interface{}
}

func Init(fs embed.FS, constants interface{}, rootTemplatePath string, commonTemplatePaths ...string) (*TemplateEngine, error) {
	root, err := template.ParseFS(fs, rootTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("error parsing root template: %w", err)
	}

	templateFunctions := template.FuncMap{
		"asdateinputvalue": AsDateInputValue,
	}

	return &TemplateEngine{constants: constants, root: template.Must(
		root.
			Funcs(templateFunctions).
			ParseFS(fs, commonTemplatePaths...),
	)}, nil
}

type TemplateEnvironment struct {
	Const interface{}
	Data  interface{}
}

func (te *TemplateEngine) Render(w http.ResponseWriter, path string, data any) error {
	tmpl := template.Must(te.root.Clone())
	tmpl, err := tmpl.ParseFiles(path)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, TemplateEnvironment{
		Const: te.constants,
		Data:  data,
	})
}
