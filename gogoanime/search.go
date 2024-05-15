package gogoanime

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Result struct {
    Title       string
    Id          string
    Url         string
    Image       string
    ReleaseDate string
    IsDub       bool
}

type SearchResult struct {
    HasNextPage bool
	CurrentPage int
    Result      []Result
}

func Search(query string, page int) SearchResult {
    baseURL := "https://anitaku.so/filter.html"
    params := url.Values{}
    params.Set("keyword", query)
    params.Set("page", string(page))

    url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

    var searchResult SearchResult

    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    searchResult.HasNextPage = doc.Find("div.anime_name.new_series > div > div > ul > li.selected").Next().Size() > 0
	searchResult.CurrentPage = page

    searchRes := doc.Find("div.last_episodes > ul > li")

    searchRes.Each(func(i int, s *goquery.Selection) {
        var item Result

        item.Title = s.Find("p.name > a").Text()
        item.Url, _ = s.Find("p.name > a").Attr("href")
        item.Image, _ = s.Find("div > a > img").Attr("src")

        releaseDate := s.Find("div.new > div.ep").Text()
        item.ReleaseDate = strings.TrimSpace(releaseDate)

        item.Id = strings.Split(item.Url, "/")[2]
        item.IsDub = strings.Contains(strings.ToLower(item.Title), "dub")

        searchResult.Result = append(searchResult.Result, item)
    })

    return searchResult
}

// data := gogoanime.Search("zero", 2)

//     fmt.Println("Has next page:", data.HasNextPage)
//     fmt.Println("Number of results:", len(data.Result))
//     for _, item := range data.Result {
//         fmt.Println("ID:", item.Id, "Title:", item.Title)
//     }
