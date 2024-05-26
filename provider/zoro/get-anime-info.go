package zoro

import (
	"bytes"
	"encoding/json"
	"github/mdtosif/go-anime/provider/interfacetype"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// var client = resty.New()

func GetAnimeInfo(animeId string)  (*interfacetype.AnimeInfo, error) {
	var animeInfo interfacetype.AnimeInfo
	// base, _ := url.Parse("https://hianime.to/")
	id := animeId[strings.LastIndex(animeId, "-")+1:]
	vidUrl := "https://hianime.to/ajax/v2/episode/list/" + id
	println(id)
	resp, err := client.R().
		//  SetHeader("Accept", "*/*").
		// SetHeader("Host", base.Host).
		SetHeader("Referer", "https://hianime.to/watch"+animeId).
		SetHeader("X-Requested-With", "XMLHttpRequest").
		Get(vidUrl)
	if err != nil {
		return nil, err
	}

	// Read the response body into a byte slice
	bodyBytes := resp.Body()


	// Define a struct to match the expected JSON structure
	type jsonResponse struct {
		Html string `json:"html"`
	}

	var response jsonResponse

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader([]byte(response.Html))
	// Parse the HTML response body
	doc, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return nil, err
	}

	doc.Find("div.detail-infor-content > div > a").Each(func(i int, s *goquery.Selection) {
		var episode interfacetype.Episode
		episode.Url = s.AttrOr("href", "")
		epId := strings.Split(episode.Url, "/")[2]
		println("epid:", epId, )
		episode.Id = epId
		animeInfo.Episodes = append(animeInfo.Episodes, episode)
	})

	return &animeInfo, nil

}
