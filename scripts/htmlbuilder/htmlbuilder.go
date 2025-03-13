package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

const (
	fileMode  = os.O_CREATE | os.O_RDWR
	filePerms = 0640
  templatePath = "../../web/base_template.html"
  pagesDir = "../../web/src/pages/"
  )

type Html struct {
		Fp    string
		Content string
	}

func resolvePath(relativePath string) (string, error) {
  _, filename, _, ok := runtime.Caller(0)
  if !ok {
    return "", fmt.Errorf("error getting source file path")
  }

  sourceDir := filepath.Dir(filename)
  return filepath.Join(sourceDir, relativePath), nil
}

func main() {
  tp, err := resolvePath(templatePath)
  if err != nil {
    panic(err)
  }
	base_template, err := os.ReadFile(tp)
	if err != nil {
		panic(err)
	}

	var bodies []Html

  dirpath, err := resolvePath(pagesDir)
  if err != nil {
    panic(err)
  }
	err = filepath.Walk(
    dirpath,
		func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			bytes, err := os.ReadFile(path)
      fp := fmt.Sprintf("../../web/build/%v", info.Name())
      fp, err = resolvePath(fp)
			bodies = append(bodies, Html{fp, string(bytes)})
			return nil
		})
	if err != nil {
		panic(err)
	}

	// Create a new template and parse the base_template into it.
	t := template.Must(template.New("base_template").Parse(string(base_template)))

	// Execute the template for each recipient.
	for _, b := range bodies {
    f, err := os.OpenFile(b.Fp, fileMode, filePerms)
    if err != nil {
      panic(err)
    }
		err = t.Execute(f, b)
		if err != nil {
			fmt.Println("error executing template:", err)
		}
	}
}
