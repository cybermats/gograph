package repository

import "time"

type tvdbInfo struct {
	Data []byte
}

type tvdbImage struct {
	Path string
}

// ImdbInfo is a Datastore entity for storing Imdb basic data
type imdbInfo struct {
	Title string // English name
	Year  int    // Year of production
}

// TitleView is a Datastore entity for storing when a title was graphed.
type titleView struct {
	TID      string    `datastore:"t_id"`     // IMDB Identifier
	Datetime time.Time `datastore:"datetime"` // Time of a graph view
}
