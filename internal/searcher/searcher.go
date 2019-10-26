package searcher

import (
	"bufio"
	"compress/gzip"
	"encoding/gob"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// EpisodeInfo hold basic info about a single episode
type EpisodeInfo struct {
	PrimaryTitle string  `json:"primary_title"`
	Season       uint8   `json:"season"`
	Episode      uint8   `json:"episode"`
	Rating       float32 `json:"rating"`
}

// TitleInfo is the representation of a single tv series title.
type TitleInfo struct {
	PrimaryTitle  string         `json:"primary_title"`
	StartYear     uint16         `json:"start_year"`
	EndYear       uint16         `json:"end_year"`
	AverageRating float32        `json:"average_rating"`
	Episodes      []*EpisodeInfo `json:"episodes"`
}

// SearchInfo is a representation of titles that is returned from search.
type SearchInfo struct {
	ID            string  `json:"t_id"`
	PrimaryTitle  string  `json:"primary_title"`
	StartYear     uint16  `json:"start_year"`
	AverageRating float32 `json:"average_rating"`
}

// Titles is a list of information for each title.
type SearchTitles []SearchInfo

func (t SearchTitles) Len() int      { return len(t) }
func (t SearchTitles) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t SearchTitles) Less(i, j int) bool {
	if t[i].AverageRating == t[j].AverageRating {
		return t[i].ID < t[j].ID
	}
	return t[i].AverageRating < t[j].AverageRating
}

// Searcher holds the database
type Searcher struct {
	titles   map[string]*TitleInfo
	episodes map[string]*EpisodeInfo
	lookup   map[string][]string
}

func onlyLettersPredicate(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c)
}

func NewSearcher() *Searcher {
	return &Searcher{
		make(map[string]*TitleInfo),
		make(map[string]*EpisodeInfo),
		make(map[string][]string),
	}
}

func NewSearcherFromFiles(basics io.Reader, episodes io.Reader, ratings io.Reader) (*Searcher, error) {
	searcher := &Searcher{
		make(map[string]*TitleInfo),
		make(map[string]*EpisodeInfo),
		make(map[string][]string),
	}
	err := readFile(basics, func(record []string) {
		id := record[0]
		titleType := filterNil(record[1])
		primary := filterNil(record[2])
		startYear := filterNil(record[5])
		endYear := filterNil(record[6])

		if titleType == "tvSeries" || titleType == "tvMiniSeries" {
			searcher.addSeason(id, titleType, primary, startYear, endYear)
		} else if titleType == "tvEpisode" {
			searcher.addEpisode(id, primary)
		}
	})
	if err != nil {
		return nil, err
	}
	err = readFile(episodes, func(record []string) {
		id := record[0]
		pID := filterNil(record[1])
		season := filterNil(record[2])
		episode := filterNil(record[3])
		searcher.linkEpisode(id, pID, season, episode)
	})
	if err != nil {
		return nil, err
	}
	err = readFile(ratings, func(record []string) {
		id := record[0]
		rating := filterNil(record[1])

		searcher.addRating(id, rating)
	})
	if err != nil {
		return nil, err
	}
	searcher.pack()

	return searcher, nil
}

func readFile(reader io.Reader, fn func(record []string)) error {
	zr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(zr)
	for scanner.Scan() {
		line := scanner.Text()
		record := strings.Split(line, "\t")
		fn(record)
	}
	return scanner.Err()
}

func filterNil(input string) string {
	if input == "\\N" {
		return ""
	}
	return input
}

func (s *Searcher) addSeason(id string, titleType string, primary string, startYear string, endYear string) {
	sy, _ := strconv.ParseInt(startYear, 10, 16)
	ey, _ := strconv.ParseInt(endYear, 10, 16)
	s.titles[id] = &TitleInfo{
		PrimaryTitle:  primary,
		StartYear:     uint16(sy),
		EndYear:       uint16(ey),
		AverageRating: 0,
		Episodes:      make([]*EpisodeInfo, 0),
	}
	words := strings.FieldsFunc(strings.ToLower(primary), onlyLettersPredicate)
	tempLookup := make(map[string]bool)

	for _, word := range words {
		if _, ok := tempLookup[word]; !ok {
			tempLookup[word] = true
			s.lookup[word] = append(s.lookup[word], id)
		}
	}
}

func (s *Searcher) addEpisode(id string, primary string) {
	eInfo := EpisodeInfo{PrimaryTitle: primary}
	s.episodes[id] = &eInfo
}

func (s *Searcher) linkEpisode(eID string, pID string, seasonNumber string, episodeNumber string) {
	title, ok := s.titles[pID]
	if !ok {
		return
	}
	season, _ := strconv.Atoi(seasonNumber)
	episode, _ := strconv.Atoi(episodeNumber)
	eInfo, ok := s.episodes[eID]
	if ok {
		eInfo.Season = uint8(season)
		eInfo.Episode = uint8(episode)
		title.Episodes = append(title.Episodes, eInfo)
	}
}

func (s *Searcher) addRating(id string, ratingStr string) {
	if episode, ok := s.episodes[id]; ok {
		rating, _ := strconv.ParseFloat(ratingStr, 32)
		episode.Rating = float32(rating)
		return
	}
	if title, ok := s.titles[id]; ok {
		rating, _ := strconv.ParseFloat(ratingStr, 32)
		title.AverageRating = float32(rating)
		return
	}
}

func (s *Searcher) pack() {
	s.episodes = nil
}

func (s *Searcher) Search(phrase string) SearchTitles {
	var superset map[string]bool
	init := false
	for _, word := range strings.FieldsFunc(strings.ToLower(phrase), onlyLettersPredicate) {
		searchSet := make(map[string]bool)
		for _, id := range s.lookup[word] {
			searchSet[id] = true
		}
		//fmt.Println("SearchSet: ", searchSet)
		if !init {
			superset = searchSet
			init = true
		} else {
			result := make(map[string]bool)
			for k := range superset {
				if _, ok := searchSet[k]; ok {
					result[k] = true
				}
			}
			superset = result
		}
		//fmt.Println("superset: ", superset)
	}
	output := make([]SearchInfo, 0, len(superset))
	for id := range superset {
		title := s.titles[id]
		output = append(output,
			SearchInfo{
				ID:            id,
				PrimaryTitle:  title.PrimaryTitle,
				StartYear:     title.StartYear,
				AverageRating: title.AverageRating,
			})
	}

	sort.Sort(sort.Reverse(SearchTitles(output)))

	return output
}

func (s *Searcher) Get(id string) (*TitleInfo, error) {
	title, ok := s.titles[id]
	if !ok {
		return nil, errors.New("Nothing found.")
	}
	return title, nil
}

func (s *Searcher) Write(writer io.Writer) error {
	data := struct {
		Titles map[string]*TitleInfo
		Lookup map[string][]string
	}{
		Titles: s.titles,
		Lookup: s.lookup,
	}
	return gob.NewEncoder(writer).Encode(data)
}

func (s *Searcher) Read(reader io.Reader) error {
	var data struct {
		Titles map[string]*TitleInfo
		Lookup map[string][]string
	}
	err := gob.NewDecoder(reader).Decode(&data)
	if err != nil {
		return err
	}
	s.titles = data.Titles
	s.lookup = data.Lookup
	return nil
}
