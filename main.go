package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Title struct {
	Id    string
	Title string
	Year  string
}

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
	topTitles := []Title{
		Title{"1", "Foo", "1951"},
		Title{"2", "Bar", "1952"},
		Title{"3", "Zoo", "1953"},
		Title{"4", "Hah", "1954"},
	}
	data := struct {
		Title  string
		Titles []Title
	}{"foo bar", topTitles}
	err := t.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func staticHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("static")
	http.Error(w, "", http.StatusInternalServerError)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *template.Template),
	tmpls *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, tmpls)
	}
}

func main() {
	tmplMap, err := readTemplates("templates", "index.html", "about.html")
	if err != nil {
		log.Fatal("Parsing templates: ", err)
		return
	}
	http.HandleFunc("/", makeHandler(mainHandler, tmplMap["index.html"]))
	http.HandleFunc("/about.html", makeHandler(aboutHandler, tmplMap["about.html"]))
	http.HandleFunc("/static/", staticHandler)
	http.HandleFunc("/favicon.ico", staticHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
