package provider

import (
	"github/mdtosif/go-anime/provider/interfacetype"
	"syscall/js"
)

type Provider interface{
 GetAnimeSrcs(episodeId string) ([]interfacetype.AnimeEpisodeSource, error)
 Search(query string, page int) (*interfacetype.SearchResult, error)
 GetAnimeInfo(animeId string) (*interfacetype.AnimeInfo, error)
}

var provider Provider

func search(this js.Value, i []js.Value) interface{} {
	arg1 := i[0].String()
	arg2 := i[0].Int()
	data, err := provider.Search(arg1, arg2)
	if err != nil {
		return err.Error()
	}
	return data
}

func getAnimeInfo(this js.Value, i []js.Value) interface{} {
	arg1 := i[0].String()
	data, err := provider.GetAnimeInfo(arg1)
	if err != nil {
		return err.Error()
	}
	return data
}

func getAnimeSrcs(this js.Value, i []js.Value) interface{} {
	arg1 := i[0].String()
	data, err := provider.GetAnimeSrcs(arg1)
	if err != nil {
		return err.Error()
	}
	return data
}

func CreateWasm(providerIn Provider){
	provider = providerIn
	js.Global().Set("search", js.FuncOf(search))
	js.Global().Set("getAnimeSrcs", js.FuncOf(getAnimeSrcs))
	js.Global().Set("getAnimeInfo", js.FuncOf(getAnimeInfo))
}

