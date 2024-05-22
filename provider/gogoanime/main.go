package gogoanime

import "github/mdtosif/go-anime/provider/interfacetype"

type GogoAnimeProvider struct{

}

func (gogo GogoAnimeProvider) GetAnimeInfo(animeId string) (*interfacetype.AnimeInfo, error){
	return getAnimeInfo(animeId)
}

func (gogo GogoAnimeProvider) GetAnimeSrcs(episodeId string)  ([]interfacetype.AnimeEpisodeSource, error){
	return getAnimeSrcs(episodeId)
}

func (gogo GogoAnimeProvider) Search(query string, page int)  (*interfacetype.SearchResult, error){
	return search(query, page)
}
