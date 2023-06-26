package pkg

import (
	"bytes"
	"fmt"
	"github.com/nestoca/actions/publish-people/go/pkg/tree"
	"html/template"
	"os"

	"github.com/nestoca/jac/pkg/config"
	"github.com/nestoca/jac/pkg/live"
)

func Render(catalogDir string, templateFile string, outputFile string) error {
	catalog, err := loadCatalog(catalogDir)
	if err != nil {
		return fmt.Errorf("loading catalog: %w", err)
	}

	tr := tree.NewTree(catalog)
	result, err := render(tr, templateFile)
	if err != nil {
		return fmt.Errorf("rendering tree: %w", err)
	}

	return os.WriteFile(outputFile, []byte(result), 0644)
}

func loadCatalog(dir string) (*live.Catalog, error) {
	cfg, err := config.LoadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	return live.LoadCatalog(cfg.Dir, cfg.Glob)
}

func render(tr *tree.Tree, templateFile string) (string, error) {
	tmpl := template.New("tmpl")
	tmpl.Funcs(template.FuncMap{
		"safeHTML": func(content string) template.HTML {
			return template.HTML(content)
		},
		"getValue": func(obj interface{}, keyPath string) (string, error) {
			valueProvider, ok := obj.(live.ValuesProvider)
			if !ok {
				return "", fmt.Errorf("object %T does not implement ValuesProvider", obj)
			}
			value, _ := valueProvider.GetValue(keyPath)
			return value, nil
		},
	})
	templateText, err := os.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("reading template file: %w", err)
	}
	tmpl, err = tmpl.Parse(string(templateText))
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", templateFile, err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, tr)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return result.String(), nil
}
