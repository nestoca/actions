package internal

import (
	"bytes"
	"fmt"
	"github.com/nestoca/jac/pkg/config"
	"github.com/nestoca/jac/pkg/live"
	"github.com/nestoca/people-renderer/internal/structure"
	"html/template"
)

func Render(catalogDir string, templateFile string) error {
	catalog, err := loadCatalog(catalogDir)
	if err != nil {
		return fmt.Errorf("loading catalog: %w", err)
	}

	tree := structure.NewTree(catalog)
	result, err := render(tree, templateFile)
	if err != nil {
		return fmt.Errorf("rendering tree: %w", err)
	}

	fmt.Println(result)

	return nil
}

func loadCatalog(dir string) (*live.Catalog, error) {
	cfg, err := config.LoadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	catalog, err := live.LoadCatalog(cfg.Dir, cfg.Glob)
	if err != nil {
		return nil, fmt.Errorf("loading catalog: %w", err)
	}
	fmt.Printf("catalog loaded from %s\n", dir)
	fmt.Printf("catalog contains %d groups\n", len(catalog.All.Groups))
	fmt.Printf("catalog contains %d people\n", len(catalog.All.People))
	return catalog, nil
}

func render(tree *structure.Tree, templateFile string) (string, error) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", templateFile, err)
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, tree)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}
	return result.String(), nil
}
