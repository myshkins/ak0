package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
  // "text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)


const postDir = "../../blog/posts"
const jsxTemplate = `
import NavBar from "../components/NavBar";
import './BlogLayout.css';
import '../index.css'

function BlogLayout() {
  return (
    <div className="blog-post">
      <NavBar />
      <h1>Blog post</h1>
    </div>
  );
}

export default BlogLayout;`

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func convertMDPostsToHTML() {
  dirPath, err := filepath.Abs(postDir)
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
    name := strings.Split(info.Name(), ".")[0]
    name = fmt.Sprintf("./html/%v.html", name)
    ap, err := filepath.Abs(name)
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println(ap)
    os.WriteFile(name, mdToHTML(bytes), 0644)
    return nil
    },
  )
}

func templateThing() {
  // execute template
  // write files to web/src/posts
}

func main() {
  /*
  convert to html
  execute template
  write output to web/src/posts/file.jsx
  */

  convertMDPostsToHTML()
	// fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, html)
}

