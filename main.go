package main

import (
	"github/mdtosif/go-anime/provider"
	"github/mdtosif/go-anime/provider/gogoanime"
)

func main()  {
	provider.CreateWasm(gogoanime.GogoAnimeProvider{})
}