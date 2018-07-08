package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kunhou/TMDB/config"
	"github.com/mlytics/micro-dns/log"

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
	dbConfig := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s sslrootcert=", cfg.Host, cfg.User, cfg.Password, cfg.DataBase, cfg.Port)
	db, err := gorm.Open("postgres", dbConfig)
	if err != nil {
		log.WithError(err).Error("Connect Database Failed")
		panic("failed to connect database")
	}
	return db
}
