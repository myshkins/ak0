//go:build !minimal

package main

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/myshkins/ak0/internal/helpers"
)

/*
take the html in web/src/pages/ (blogger output)
execute template to add navbar and css link
output to web/build/
*/

const (
	fileMode              = os.O_CREATE | os.O_RDWR
	dirPerms              = 0777
	filePerms             = 0666
	baseTemplatePath      = "../../web/base_template.html"
	blogIndexTemplatePath = "../../web/blog_index_template.html"
	indexPath             = "../../web/src/index.html"
	blogPath              = "../../web/src/blog/index.html"
	pagesDir              = "../../web/src/blog/posts/"
	blogOutDir            = "../../web/build/blog/"
	outDir                = "../../web/build/"
	indexCss              = "/assets/index.css"
	blogIndexCss          = "/assets/blogIndex.css"
	homeCss               = "/assets/home.css"
	postCss               = "/assets/post.css"
)

// todo: add `Updated` field. would have to check if post file exists in git?
type PostContent struct {
	Fp         string
	Styles     []string
	IsBlogPost bool
	Content    string
}

type BlogIndexContent struct {
	Posts      []string
	Styles     []string
	IsBlogPost bool
	Content    string
}

func main() {
	tp, err := helpers.MakeRelPathAbs(baseTemplatePath)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
	baseTemplateFile, err := os.ReadFile(tp)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

  btp, err := helpers.MakeRelPathAbs(blogIndexTemplatePath)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
	blogTemplateFile, err := os.ReadFile(btp)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

	var bodies []PostContent

  biContent := BlogIndexContent{Styles: []string{indexCss, blogIndexCss}}

	dirpath, err := helpers.MakeRelPathAbs(pagesDir)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
	absOutdir, err := helpers.MakeRelPathAbs(outDir)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

	absBlgOutDir, err := helpers.MakeRelPathAbs(blogOutDir)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
	err = os.MkdirAll(absBlgOutDir, dirPerms)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

	// walk src/blog/posts dir, create postContent with fp and content, add to bodies[]
	// also add title of each post to biContent
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
			fmt.Printf("post name is: %s\n", name)
			biContent.Posts = append(biContent.Posts, name)
			postOutDir, err := helpers.MakeRelPathAbs(filepath.Join(blogOutDir, name))
			if err != nil {
				slog.Error("fatal error", "error", err)
				os.Exit(1)
			}
			err = os.Mkdir(postOutDir, dirPerms)
			bytes, err := os.ReadFile(path)
			fp := filepath.Join(postOutDir, "index.html")
			bodies = append(
				bodies,
				PostContent{
					Fp:         fp,
					Styles:     []string{postCss, indexCss},
					IsBlogPost: true,
					Content:    string(bytes),
				},
			)
			return nil
		})
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

	// handle index.html
	ifp, err := helpers.MakeRelPathAbs(indexPath)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}
	bytes, err := os.ReadFile(ifp)
	fp := filepath.Join(absOutdir, "index.html")
	bodies = append(
		bodies,
		PostContent{
			Fp:         fp,
			Styles:     []string{indexCss, homeCss},
			IsBlogPost: false,
			Content:    string(bytes),
		},
	)

	// handle blogIndex.html
	// for each post add a new li item to index with title and date of post
	// Create a blog index template and parse the blogIndex into it.
	blogT := template.Must(template.New("blog-index").Parse(string(blogTemplateFile)))
  blogIndexFp := filepath.Join(absBlgOutDir, "index.html")
	f, err := os.OpenFile(blogIndexFp, fileMode, filePerms)
	if err != nil {
		slog.Error("fatal error", "error", err)
		os.Exit(1)
	}

	err = blogT.Execute(f, biContent)
	if err != nil {
		fmt.Println("error executing template:", err)
	}

	// Create base template and parse the base_template into it.
	t := template.Must(template.New("base").Parse(string(baseTemplateFile)))

	// Execute the template for each page body
	for _, b := range bodies {
		f, err := os.OpenFile(b.Fp, fileMode, filePerms)
		if err != nil {
			slog.Error("fatal error", "error", err)
			os.Exit(1)
		}
		err = t.Execute(f, b)
		if err != nil {
			fmt.Println("error executing template:", err)
		}
	}
}
