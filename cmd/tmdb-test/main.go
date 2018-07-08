package main

import (
	"fmt"

	"github.com/ryanbradynd05/go-tmdb"
)

func main() {
	client := tmdb.Init("df62c2ec98e3366767d418bae01180fa")
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
