package main

import (
	"flag"

	"github.com/kunhou/TMDB/db"
	"github.com/kunhou/TMDB/router"
)

func main() {
	// migration
	var drop = flag.Bool("drop", false, "drop all tables")
	var rollback = flag.Int("rollback", 0, "rollback how many steps")
	flag.Parse()
	if *drop {
		db.Drop()
	}

	db.Migrate(*rollback)
	router.Setting().Run()
}
