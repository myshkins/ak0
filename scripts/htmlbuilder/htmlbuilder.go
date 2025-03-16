//go:build !minimal

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

  "github.com/alecthomas/chroma/v2"

	"github.com/myshkins/ak0/internal/helpers"
)

/*
take the html in web/src/pages/ (blogger output)
execute template to add navbar and css link
output to web/build/
*/

const (
	fileMode  = os.O_CREATE | os.O_RDWR
  dirPerms = 0777
	filePerms = 0666
  templatePath = "../../web/base_template.html"
  indexPath = "../../web/src/index.html"
  blogPath = "../../web/src/blogIndex.html"
  pagesDir = "../../web/src/pages/"
  blogOutDir = "../../web/build/blog/"
  outDir = "../../web/build/"
  indexCss = "/assets/index.css"
  homeCss = "/assets/home.css"
  postCss = "/assets/post.css"
  )

type Html struct {
		Fp    string
    Styles  []string
    IsBlogPost bool
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

  // need to create necessary directories
  // web/build/
  // web/build/blog/
  // web/build/blog/eachpost/
  // use file system as url naviagtion
  // each dir contains an index.html
  absBlgOutDir, err := helpers.ResolvePath(blogOutDir)
  if err != nil {
    panic(err)
  }
  err = os.MkdirAll(absBlgOutDir, dirPerms)
  if err != nil {
    panic(err)
  }
  // gotta fix walkdir to write to separate dirs

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
      name := strings.Split(info.Name(), ".")[0]
      postOutDir, err := helpers.ResolvePath(filepath.Join(blogOutDir, name))
      if err != nil {
        panic(err)
      }
      err = os.Mkdir(postOutDir, dirPerms)
			bytes, err := os.ReadFile(path)
      fp := filepath.Join(postOutDir, "index.html")
			bodies = append(
        bodies,
        Html{
          fp,
          []string{postCss, indexCss},
          true,
          string(bytes),
        },
      )
			return nil
		})
	if err != nil {
		panic(err)
	}

  // handle index.html
  ifp, err := helpers.ResolvePath(indexPath)
  if err != nil {
    panic(err)
  }
  bytes, err := os.ReadFile(ifp)
  fp := filepath.Join(absOutdir, "index.html")
  bodies = append(
    bodies,
    Html{
      fp,
      []string{indexCss, homeCss},
      false,
      string(bytes),
    },
  )

  // handle blogIndex.html
  bfp, err := helpers.ResolvePath(blogPath)
  if err != nil {
    panic(err)
  }
  bytes, err = os.ReadFile(bfp)
  fp = filepath.Join(absBlgOutDir, "index.html")
  bodies = append(
    bodies,
    Html{
      fp,
      []string{"/assets/index.css", "/assets/blogIndex.css"},
      false,
      string(bytes),
    },
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
