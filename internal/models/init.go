package models

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/CheerChen/konachan-app/internal/cache"
	"github.com/CheerChen/konachan-app/internal/log"
)

var db *gorm.DB
var cc *cache.Handler

func Init() {
	db = GetDB()
	cc = &cache.Handler{"./cache", 60}
}

func GetDB() *gorm.DB {
	dbFile := "sqlite.db"
	if _, err := os.Stat(dbFile); err != nil {
		_, err = os.Create(dbFile)
		if err != nil {
			log.Fatalf("Init DB failed: %s", err)
		}
	}

	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Init DB failed: %s", err)
	}

	// 自动迁移模式
	err = db.AutoMigrate(&Post{}).Error
	if err != nil {
		log.Fatalf("AutoMigrate failed: %s", err)
	}
	return db
}
