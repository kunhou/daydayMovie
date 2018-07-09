package main

import (
	"fmt"

	"github.com/kunhou/TMDB/config"
	"github.com/ryanbradynd05/go-tmdb"
)

var cfg = config.GetConfig()

func main() {
	client := tmdb.Init(cfg.TMDBToken)
	var options = make(map[string]string)
	options["language"] = "zh-TW"
	fightClubInfo, err := client.DiscoverMovie(options)
	// fightClubInfo, err := client.GetMovieInfo(348350, options)

	if err != nil {
		fmt.Println(err)
	}
	fightClubJSON, err := tmdb.ToJSON(fightClubInfo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", fightClubJSON)
}
