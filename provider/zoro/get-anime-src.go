package zoro

import (
	"bytes"
	"encoding/json"
	"github/mdtosif/go-anime/provider/interfacetype"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// var client = resty.New()

func GetAnimeSrcId(animeId string)  ([]string, error) {
	id := strings.Split(animeId, "?ep=")[1]
	// base, _ := url.Parse("https://hianime.to/")
	vidUrl := "https://hianime.to/ajax/v2/episode/servers?episodeId=" + id
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

	var srcIds []string
	doc.Find("div.servers-sub div.item").Each(func(i int, s *goquery.Selection) {
		srcId := s.AttrOr("data-id", "")
		println("srcId: ", srcId)

		srcIds = append(srcIds, srcId)
	})

	return srcIds, nil

}


func GetAnimeSrcs(animeId string)   (*interfacetype.AnimeInfo, error){
	var animeInfo interfacetype.AnimeInfo
	id, err := GetAnimeSrcId(animeId)
	if err != nil {
		return nil, err
	}
	// base, _ := url.Parse("https://hianime.to/")
	vidUrl := "https://hianime.to/ajax/v2/episode/sources?id=" + id[1]
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
	println(string(bodyBytes))
	
	// Define a struct to match the expected JSON structure
	type jsonResponse struct {
		Type string `json:"type"`
		Link string `json:"link"`
	}

	var response jsonResponse

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}

	

	return &animeInfo, nil
}