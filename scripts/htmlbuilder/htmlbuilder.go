//go:build !minimal
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"

  "github.com/myshkins/ak0/internal/helpers"
)

/*
take the html in web/src/pages/ (blogger output)
execute template to add navbar and css link
output to web/build/
*/

const (
	fileMode  = os.O_CREATE | os.O_RDWR
	filePerms = 0640
  templatePath = "../../web/base_template.html"
  indexPath = "../../web/src/index.html"
  pagesDir = "../../web/src/pages/"
  outDir = "../../web/build/"
  )

type Html struct {
		Fp    string
    Styles  []string
		Content string
	}

func main() {
  tp, err := helpers.ResolvePath(templatePath)
  if err != nil {
    panic(err)
  }
	base_template, err := os.ReadFile(tp)
	if err != nil {
		panic(err)
	}

	var bodies []Html

  dirpath, err := helpers.ResolvePath(pagesDir)
  if err != nil {
    panic(err)
  }
  absOutdir, err := helpers.ResolvePath(outDir)
    if err != nil {
      panic(err)
    }

  // walk src/pages dir, add fp and content to bodies[]
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
      fp := filepath.Join(absOutdir, info.Name())
			bodies = append(
        bodies,
        Html{fp, []string{"assets/post.css", "assets/index.css"}, string(bytes)},
      )
			return nil
		})
	if err != nil {
		panic(err)
	}

  // add logic here to handle index. just need to add it to the body slice?
  ifp, err := helpers.ResolvePath(indexPath)
  if err != nil {
    panic(err)
  }
  bytes, err := os.ReadFile(ifp)
  fp := filepath.Join(absOutdir, "index.html")
  bodies = append(
    bodies,
    Html{fp, []string{"assets/index.css", "assets/home.css"}, string(bytes)},
  )



	// Create a new template and parse the base_template into it.
	t := template.Must(template.New("base_template").Parse(string(base_template)))

	// Execute the template for each page body
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
