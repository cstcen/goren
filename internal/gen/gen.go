package gen

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates
var content embed.FS

type Params struct {
	ModuleName        string
	Version           string
	PackageNameApi    string
	PackageNameApiGen string
}

func Generate(outputFile string, ops Params) error {
	t := template.New("goren")
	if err := LoadTemplates(content, t); err != nil {
		return err
	}
	if err := GenerateConfig(t, outputFile, ops); err != nil {
		return err
	}
	if err := GenerateConfigSchemas(t, outputFile, ops); err != nil {
		return err
	}
	if err := GenerateConfigParameters(t, outputFile, ops); err != nil {
		return err
	}
	if err := GenerateConfigResponses(t, outputFile, ops); err != nil {
		return err
	}
	if err := GenerateMain(t, ops); err != nil {
		return err
	}
	return GenerateGoren(t, outputFile, ops)
}

func GenerateConfig(t *template.Template, file string, ops Params) error {
	configOut, err := GenerateTemplates([]string{"config.tmpl"}, t, ops)
	if err != nil {
		return err
	}
	dir := filepath.Dir(file)
	realPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	prevDir := filepath.Join(realPath, "./")
	filename := filepath.Join(prevDir, "goren-config.yaml")
	if _, err := os.ReadFile(filename); os.IsNotExist(err) {
		if err := os.WriteFile(filename, []byte(configOut), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GenerateConfigParameters(t *template.Template, file string, ops Params) error {
	configOut, err := GenerateTemplates([]string{"config-parameters.tmpl"}, t, ops)
	if err != nil {
		return err
	}
	dir := filepath.Dir(file)
	realPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	outputDir := filepath.Join(realPath, "./parameters")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			os.Exit(-1)
		}
	}
	filename := filepath.Join(outputDir, "goren-config-parameters.yaml")
	if _, err := os.ReadFile(filename); os.IsNotExist(err) {
		if err := os.WriteFile(filename, []byte(configOut), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GenerateConfigResponses(t *template.Template, file string, ops Params) error {
	configOut, err := GenerateTemplates([]string{"config-responses.tmpl"}, t, ops)
	if err != nil {
		return err
	}
	dir := filepath.Dir(file)
	realPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	outputDir := filepath.Join(realPath, "./responses")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			os.Exit(-1)
		}
	}
	filename := filepath.Join(outputDir, "goren-config-responses.yaml")
	if _, err := os.ReadFile(filename); os.IsNotExist(err) {
		if err := os.WriteFile(filename, []byte(configOut), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GenerateConfigSchemas(t *template.Template, file string, ops Params) error {
	configOut, err := GenerateTemplates([]string{"config-schemas.tmpl"}, t, ops)
	if err != nil {
		return err
	}
	dir := filepath.Dir(file)
	realPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	outputDir := filepath.Join(realPath, "./schemas")
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			os.Exit(-1)
		}
	}
	filename := filepath.Join(outputDir, "goren-config-schemas.yaml")
	if _, err := os.ReadFile(filename); os.IsNotExist(err) {
		if err := os.WriteFile(filename, []byte(configOut), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GenerateGoren(t *template.Template, file string, ops Params) error {
	gorenOut, err := GenerateTemplates([]string{"goren.tmpl"}, t, ops)
	if err != nil {
		return err
	}
	dir := filepath.Dir(file)
	realPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	prevDir := filepath.Join(realPath, "../")
	filename := filepath.Join(prevDir, "goren.go")
	if _, err := os.ReadFile(filename); os.IsNotExist(err) {
		if err := os.WriteFile(filename, []byte(gorenOut), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func GenerateMain(t *template.Template, ops Params) error {
	mainOut, err := GenerateTemplates([]string{"main.tmpl"}, t, ops)
	if err != nil {
		return err
	}

	realPath, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	filename := filepath.Join(realPath, "main.go")
	if _, err := os.ReadFile(filename); os.IsNotExist(err) {
		if err := os.WriteFile(filename, []byte(mainOut), os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

// GenerateTemplates used to generate templates
func GenerateTemplates(templates []string, t *template.Template, ops interface{}) (string, error) {
	var generatedTemplates []string
	for _, tmpl := range templates {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		if err := t.ExecuteTemplate(w, tmpl, ops); err != nil {
			return "", fmt.Errorf("error generating %s: %s", tmpl, err)
		}
		if err := w.Flush(); err != nil {
			return "", fmt.Errorf("error flushing output buffer for %s: %s", tmpl, err)
		}
		generatedTemplates = append(generatedTemplates, buf.String())
	}

	return strings.Join(generatedTemplates, "\n"), nil
}

// LoadTemplates loads all of our template files into a text/template. The
// path of template is relative to the templates' directory.
func LoadTemplates(src embed.FS, t *template.Template) error {
	return fs.WalkDir(src, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %s: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}

		buf, err := src.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file '%s': %w", path, err)
		}

		templateName := strings.TrimPrefix(path, "templates/")
		tmpl := t.New(templateName)
		_, err = tmpl.Parse(string(buf))
		if err != nil {
			return fmt.Errorf("parsing template '%s': %w", path, err)
		}
		return nil
	})
}
