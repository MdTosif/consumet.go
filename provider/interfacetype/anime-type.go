package interfacetype

type Episode struct {
	Id     string
	Number string
	Url    string
}

type AnimeInfo struct {
	Title        string
	Image        string
	ReleaseDate  string
	IsDub        bool
	Type         string
	Status       string
	Description  string
	TotalEpisode int
	Url          string
	Episodes     []Episode
}

type AnimeEpisodeSource struct {
	Quality string
	Url     string
	Name    string
}

type Result struct {
	Title       string `json:"title"`
	Id          string `json:"id"`
	Url         string `json:"url"`
	Image       string `json:"image"`
	ReleaseDate string `json:"releaseDate"`
	IsDub       bool   `json:"isDub"`
}

type SearchResult struct {
	HasNextPage bool `json:"hasNextPage"`
	CurrentPage int  `json:"currentPage"`
	Result      []Result
}
