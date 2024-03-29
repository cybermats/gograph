package repository

// TitleTopInfo contains information about how much a title has been viewed in the last period of time
type TitleTopInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
	Count int    `json:"count"`
}
