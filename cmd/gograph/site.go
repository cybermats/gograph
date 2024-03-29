package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"cybermats/gograph/internal/repository"
	"cybermats/gograph/internal/searcher"

	"github.com/gorilla/mux"
)

func aboutHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	log.Println("about")
	err := t.ExecuteTemplate(w, "about", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request, t *template.Template) {
	log.Println("main")
	topTitles, err := repository.GetTop(7, 3)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		Title  string
		Titles []repository.TitleTopInfo
	}{"", topTitles}

	err = t.ExecuteTemplate(w, "index", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type rating struct {
	X    []float32
	Y    []float32
	Text []string
	Name string
}

func extractGraphsAndTrendlines(episodes *searcher.EpisodesInfo) ([]rating, trendline, []trendline) {
	sort.Sort(episodes)

	ratings := make([]rating, 0)
	season := uint8(255)
	seasonIdx := 0
	xaxis := make([]float32, 0)
	yaxis := make([]float32, 0)
	text := make([]string, 0)

	trendlines := make([]trendline, 0)

	showMaker := trendlineMaker{}
	seasonMaker := trendlineMaker{}

	for idx, episode := range *episodes {
		newSeason := season != episode.Season
		if newSeason && season != 255 {
			ratings = append(ratings, rating{
				X:    xaxis,
				Y:    yaxis,
				Text: text,
				Name: fmt.Sprintf("Season %d", season),
			})
			trendlines = append(trendlines, seasonMaker.trendline(seasonIdx))

			xaxis = make([]float32, 0, len(xaxis))
			yaxis = make([]float32, 0, len(yaxis))
			text = make([]string, 0, len(text))
			seasonMaker = trendlineMaker{}
			seasonIdx = idx
		}

		x := float32(idx)
		y := episode.Rating
		showMaker.addSample(x, y)
		seasonMaker.addSample(x, y)
		xaxis = append(xaxis, x)
		yaxis = append(yaxis, y)
		text = append(text,
			fmt.Sprintf("S%dE%d - %s", episode.Season,
				episode.Episode, episode.PrimaryTitle))

		season = episode.Season
	}
	ratings = append(ratings, rating{
		X:    xaxis,
		Y:    yaxis,
		Text: text,
		Name: fmt.Sprintf("Season %d", season),
	})
	trendlines = append(trendlines, seasonMaker.trendline(seasonIdx))

	return ratings, showMaker.trendline(0), trendlines
}

func graphHandler(w http.ResponseWriter, r *http.Request, t *template.Template, db *searcher.Db) {
	log.Println("graph handler")
	id := mux.Vars(r)["id"]
	result, err := db.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ratings, showTrendline, seasonTrendline := extractGraphsAndTrendlines(&result.Episodes)
	tvdbInfo, tvdbImage, err := getTVDBMetaData(id)

	data := struct {
		ID            string
		Title         string
		Titles        []string
		Year          uint16
		AverageRating float32
		Ratings       []rating
		Trendline     trendline
		Trendlines    []trendline
		Info          interface{}
		Image         string
	}{
		ID:            id,
		Title:         result.PrimaryTitle,
		Titles:        make([]string, 0),
		Year:          result.StartYear,
		AverageRating: result.AverageRating,
		Ratings:       ratings,
		Trendline:     showTrendline,
		Trendlines:    seasonTrendline,
		Info:          tvdbInfo,
		Image:         tvdbImage,
	}

	err = t.ExecuteTemplate(w, "graph", data)
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

func makeTemplateAndDatabaseHandler(
	fn func(http.ResponseWriter, *http.Request, *template.Template, *searcher.Db),
	tmpls *template.Template, db *searcher.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, tmpls, db)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func initSite(router *mux.Router, webDirectory string, db *searcher.Db) error {
	s := router.PathPrefix("/").Subrouter()

	tmplMap, err := readTemplates(
		filepath.Join(webDirectory, "templates"),
		"index.html", "about.html", "graph.html")
	if err != nil {
		return err
	}
	s.HandleFunc("/",
		makeTemplateHandler(mainHandler, tmplMap["index.html"]))
	s.HandleFunc("/{id:tt[0-9]+}",
		makeTemplateAndDatabaseHandler(graphHandler, tmplMap["graph.html"], db))
	s.HandleFunc("/about.html",
		makeTemplateHandler(aboutHandler, tmplMap["about.html"]))

	dir := filepath.Join(webDirectory, "static")

	fs := http.FileServer(http.Dir(dir))
	s.PathPrefix("/static").Handler(http.StripPrefix("/static/", fs))
	s.Handle("/favicon.ico", fs)

	s.Use(loggingMiddleware)

	return nil
}
