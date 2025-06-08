package internal

import (
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

func LoadSpecFromFile(loader *openapi3.Loader, specPath string) (*openapi3.T, error) {
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("loader.LoadFromFile: %w", err)
	}
	return doc, nil
}

func WriteSpecToFile(doc *openapi3.T, specPath string) error {
	yamlData, err := doc.MarshalYAML()
	if err != nil {
		return fmt.Errorf("doc.MarshalYAML: %w", err)
	}

	outputFile, err := os.Create(specPath)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer outputFile.Close() //nolint:errcheck

	encoder := yaml.NewEncoder(outputFile)
	encoder.SetIndent(2)
	defer encoder.Close() //nolint:errcheck

	if err := encoder.Encode(yamlData); err != nil {
		return fmt.Errorf("encoder.Encode: %w", err)
	}
	return nil
}
