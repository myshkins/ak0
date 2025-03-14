package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/myshkins/ak0/internal/helpers"
)

const (
  postDir = "../../blog/posts"
  outDir = "../../web/src/pages/"
)


func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func main() {
  dirPath, err := helpers.ResolvePath(postDir)
  if err != nil {
    fmt.Println(err.Error())
  }
  filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
    if err != nil {
      fmt.Println(err)
      return err
    }
    if info.IsDir() {
      fmt.Printf("skipping dir: %v\n", info.Name())
      return nil
    } 
    bytes, err := os.ReadFile(path)
    // add .html extension
    name := strings.Split(info.Name(), ".")[0]
    dir, err := helpers.ResolvePath(outDir)
    if err != nil {
      panic(err)
    }
    fmt.Printf("outDir = %v\n", dir)
    fp := filepath.Join(dir, name)
    os.WriteFile(fp, mdToHTML(bytes), 0644)
    return nil
    },
  )
}

