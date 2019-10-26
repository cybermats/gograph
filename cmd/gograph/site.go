package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"cybermats/gograph/internal/repository"

	"github.com/gorilla/mux"
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
		Titles []repository.TitleTopInfo
	}{"foo bar", topTitles}

	err = t.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func graphHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	log.Println("graph handler")
}

func makeTemplateHandler(
	fn func(http.ResponseWriter, *http.Request, *template.Template),
	tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, tmpls)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func initSite(router *mux.Router, webDirectory string) error {
	s := router.PathPrefix("/").Subrouter()

	tmplMap, err := readTemplates(
		filepath.Join(webDirectory, "templates"),
		"index.html", "about.html")
	if err != nil {
		return err
	}
	s.HandleFunc("/",
		makeTemplateHandler(mainHandler, tmplMap["index.html"]))
	s.HandleFunc("/{id:tt[0-9]+}",
		makeTemplateHandler(graphHandler, tmplMap["index.html"]))
	s.HandleFunc("/about.html",
		makeTemplateHandler(aboutHandler, tmplMap["about.html"]))

	dir := filepath.Join(webDirectory, "static")

	fs := http.FileServer(http.Dir(dir))
	s.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))
	s.Handle("/favicon.ico", fs)

	s.Use(loggingMiddleware)

	return nil
}
