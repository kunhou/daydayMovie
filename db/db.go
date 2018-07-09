package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kunhou/TMDB/config"
	"github.com/kunhou/TMDB/log"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB
var cfg = config.GetConfig().DB

func init() {
	DB = NewGorm()
	DB.DB().SetConnMaxLifetime(1 * time.Minute)
	DB.DB().SetMaxOpenConns(10)
	DB.DB().SetMaxIdleConns(0)
}

// NewGorm creates a gorm DB instance
func NewGorm() *gorm.DB {
	connStr := "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable sslrootcert="
	if cfg.SSLEnable {
		connStr += fmt.Sprintf("sslmode=require sslrootcert=%s", cfg.SSLPath)
	} else {
		connStr += "sslmode=disable"
	}
	dbConfig := fmt.Sprintf(connStr, cfg.Host, cfg.User, cfg.Password, cfg.DataBase, cfg.Port)
	db, err := gorm.Open("postgres", dbConfig)
	if err != nil {
		log.WithError(err).Error("Connect Database Failed")
		panic("failed to connect database")
	}
	return db
}
