//go:build !minimal

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	chroma "github.com/alecthomas/chroma/v2"
	chromaHtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/myshkins/ak0/internal/helpers"
)

/*
read md post file in blog/posts
conver to html
  if code block:
    use a node render hook and render with chroma
write to web/src/blog/posts/
*/

const (
  dirPerms = 0777
  postDir = "../../blog/posts"
  outDir = "../../web/src/blog/posts/"
  codeSyntaxStyle = "catppuccin-mocha"
)


func renderCodeWithChroma(source, language string) (string, error) {
  var buf bytes.Buffer
  lexer := lexers.Get(language)
	if lexer == nil {
    fmt.Println("failed to get specified lexer")
		lexer = lexers.Analyse(source)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
  lexer = chroma.Coalesce(lexer)

  style := styles.Get(codeSyntaxStyle)
	if style == nil {
		style = styles.Fallback
	}

  formatter := chromaHtml.New(chromaHtml.WithClasses(true))
  
  iterator, err := lexer.Tokenise(nil, source)
  if err != nil {
    return "", err
	}

  formatter.Format(&buf, style, iterator)
  _, err = buf.Write([]byte("<style>"))
  if err != nil {
    return "", err
  }
  err = formatter.WriteCSS(&buf, style)
  if err != nil {
    return "", err
  }
  _, err = buf.Write([]byte("</style>"))
  if err != nil {
    return "", err
  }

  return buf.String(), nil
}

func renderHookCodeBlock(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
  codeNode, ok := node.(*ast.CodeBlock)
  if !ok {
    return ast.GoToNext, false
  }
  renderedCode, err := renderCodeWithChroma(
    string(codeNode.Literal),
    string(codeNode.Info),
    )
  if err != nil {
    panic(err)
  }
  io.WriteString(w, renderedCode)
  return ast.GoToNext, true
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.NoEmptyLineBeforeBlock | parser.Attributes
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.LazyLoadImages
	opts := html.RendererOptions{
    Flags: htmlFlags,
    RenderNodeHook: renderHookCodeBlock,
  }
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func main() {
  absPostDirPath, err := helpers.ResolvePath(postDir)
  if err != nil {
    panic(err)
  }
  absOutDir, err := helpers.ResolvePath(outDir)
  if err != nil {
    panic(err)
  }
  err = os.MkdirAll(absOutDir, dirPerms)
  if err != nil {
    panic(err)
  }
  filepath.Walk(absPostDirPath, func(path string, info fs.FileInfo, err error) error {
    if err != nil {
      fmt.Println(err)
      return err
    }
    if info.IsDir() {
      return nil
    } 
    bytes, err := os.ReadFile(path)
    // add .html extension
    name := strings.Split(info.Name(), ".")[0]
    name = strings.Join([]string{name, ".html"}, "")
    fp := filepath.Join(absOutDir, name)
    err = os.WriteFile(fp, mdToHTML(bytes), 0644)
    if err != nil {
      fmt.Println(err.Error())
    }
    return nil
    },
  )
}

