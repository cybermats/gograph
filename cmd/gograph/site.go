package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"cybermats/gograph/internal/repository"
)

func readTemplates(dir string, files ...string) (map[string]*template.Template, error) {
	funcs := template.FuncMap{"add": func(a int, b int) int { return a + b }}
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

func aboutHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	log.Println("about")
	err := t.ExecuteTemplate(w, "about", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	//	id := r.URL.Path[1:]

	log.Println("main")

	topTitles, err := repository.GetTop(7, 3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		Title  string
		Titles []repository.TitleInfo
	}{"foo bar", topTitles}

	err = t.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeTemplateHandler(
	fn func(http.ResponseWriter, *http.Request, *template.Template),
	tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, tmpls)
	}
}

func initSite() error {
	tmplMap, err := readTemplates("templates", "index.html", "about.html")
	if err != nil {
		return err
	}
	http.HandleFunc("/",
		makeTemplateHandler(mainHandler, tmplMap["index.html"]))
	http.HandleFunc("/about.html",
		makeTemplateHandler(aboutHandler, tmplMap["about.html"]))

	return nil
}
