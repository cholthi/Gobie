package utils

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/**
* Use it to parse nested templates in a given dir. Unlike template.ParseFiles, it tranvseres
* recursively in to the dir tree and allows parsing of all templates.
*
*you can use like this `t, err := ParseTemplates("templates",funcs)`
*  t.ExecuteTemplate(w,"view/index.html")
 */
func ParseTemplates(rootDir string, funcs template.FuncMap) (*template.Template, error) {
	cleanedRootDir := filepath.Clean(rootDir)
	pfx := len(cleanedRootDir) + 1
	rootTemplate := template.New("")

	err := filepath.Walk(cleanedRootDir, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}
			templateName := path[pfx:]
			t := rootTemplate.New(templateName).Funcs(funcs)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}
		}
		return nil
	})
	return rootTemplate, err
}
