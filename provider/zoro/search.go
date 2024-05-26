package zoro

import (
	"bytes"
	"fmt"
	"github/mdtosif/go-anime/provider/interfacetype"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func Search(keyword string, page int64) (*interfacetype.SearchResult, error) {
	var searchResult interfacetype.SearchResult
	vidUrl := "https://hianime.to/search?keyword=" + keyword + "&page=" +  fmt.Sprintf("%d", page)
	resp, err := client.R().
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Get(vidUrl)
	if err != nil {
		return nil, err
	}

	// Read the response body into a byte slice
	bodyBytes := resp.Body()
	bodyReader := bytes.NewReader(bodyBytes)

	// Parse the HTML response body
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}

	// Find and print the text of the elements matching the selector
	doc.Find("div.film-detail a").Each(func(index int, elem *goquery.Selection) {
		var item interfacetype.Result
		item.Title = elem.Text()
		item.Id = strings.Split( elem.AttrOr("href", ""), "?ref=search")[0]
		searchResult.Result = append(searchResult.Result, item)
	})

	return &searchResult, nil
}
