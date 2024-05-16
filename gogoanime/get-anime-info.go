package gogoanime

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
	Episodes     []Episode
}

func GetAnimeInfo(animeId string) AnimeInfo {
	// id fairy-tail-2014
	baseURL := "https://anitaku.so/category/" + animeId
	var animeInfo AnimeInfo

	resp, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	animeInfo.Title = strings.TrimSpace(doc.Find("section.content_left > div.main_body > div:nth-child(2) > div.anime_info_body_bg > h1").Text())
	animeInfo.Image = doc.Find("div.anime_info_body_bg > img").AttrOr("src", "")
	animeInfo.ReleaseDate = strings.Split(strings.TrimSpace(doc.Find("div.anime_info_body_bg > p:nth-child(8)").Text()), "Released: ")[1]
	animeInfo.Description = strings.Split(strings.TrimSpace(doc.Find("div.anime_info_body_bg > div:nth-child(6)").Text()), "Plot Summary: ")[0]
	animeInfo.IsDub = strings.Contains(strings.ToLower(animeInfo.Title), "dub")
	animeInfo.Type = strings.ToUpper(strings.TrimSpace(doc.Find("div.anime_info_body_bg > p:nth-child(4) > a").Text()))
	animeInfo.Status = strings.ToUpper(strings.TrimSpace(doc.Find("div.anime_info_body_bg > p:nth-child(9) > a").Text()))

	// fetching episodes
	epStart, _ := doc.Find("#episode_page > li").First().Find("a").Attr("ep_start")
	epEnd, _ := doc.Find("#episode_page > li").Last().Find("a").Attr("ep_end")
	movieId, _ := doc.Find("#movie_id").Attr("value")
	alias, _ := doc.Find("#alias_anime").Attr("value")

	te, _ := strconv.ParseInt(epEnd, 10, 0)
	animeInfo.TotalEpisode = int(te)

	params := url.Values{}
	params.Set("ep_start", epStart)
	params.Set("ep_end", epEnd)
	params.Set("id", movieId)
	params.Set("default_ep", "0")
	params.Set("alias", alias)

	epList, err := http.Get("https://ajax.gogocdn.net/ajax/load-list-episode?" + params.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer epList.Body.Close()

	doc2, err := goquery.NewDocumentFromReader(epList.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc2.Find("#episode_related > li").Each(func(i int, s *goquery.Selection) {
		var episode Episode
		href := s.Find("a").AttrOr("href", "")
		episode.Id = strings.Split(href, "/")[1]
		episode.Number = strings.Replace(s.Find("div.name").Text(), "EP ", "", -1)
		episode.Url = strings.TrimSpace(href)
		animeInfo.Episodes = append(animeInfo.Episodes, episode)
	})
	return animeInfo
}

// data:=gogoanime.GetAnimeInfo("fairy-tail-2014")

// 	for _,v := range data.Episodes{
// 		println(v.Id, v.Url)
// 	}
