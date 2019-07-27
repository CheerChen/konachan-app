package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/CheerChen/konachan-app/internal/cache"
	"github.com/CheerChen/konachan-app/internal/log"
)

var db *gorm.DB
var cc *cache.Handler

func Init() {
	GetDB()
	cc = &cache.Handler{"./cache", 60}
}

func GetDB() {
	db, err := gorm.Open("sqlite3", "db/sqlite.db")
	if err != nil {
		log.Fatalf("Init DB failed: %s", err)
	}
	defer db.Close()

	// 自动迁移模式
	db.AutoMigrate(&Post{})
}
