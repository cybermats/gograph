package main

import (
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
)

func readTemplates(dir string, files ...string) (map[string]*template.Template, error) {
	funcs := template.FuncMap{
		"add":  func(a int, b int) int { return a + b },
		"mod":  func(a int, b int) int { return a % b },
		"join": strings.Join,
		"makeRange": func(start, end int) []int {
			o := make([]int, end-start+1)
			for i := range o {
				o[i] = start + i
			}
			return o
		},
		"intJoin": func(a []int, d string) string {
			o := make([]string, len(a))
			for i, v := range a {
				o[i] = strconv.Itoa(v)
			}
			return strings.Join(o, d)
		},
	}
	pattern := filepath.Join(dir, "helpers", "*.html")
	baseTemplates := template.Must(template.New("root").Funcs(funcs).ParseGlob(pattern))
	tmplMap := make(map[string]*template.Template)
	for _, file := range files {
		tmpl, err := baseTemplates.Clone()
		if err != nil {
			return nil, err
		}
		pattern = filepath.Join(dir, file)
		_, err = tmpl.ParseFiles(pattern)
		if err != nil {
			return nil, err
		}
		tmplMap[file] = tmpl
	}

	return tmplMap, nil
}
