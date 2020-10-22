package dal

import (
	"fmt"
	"log"

	"nicheanal.com/config"

	"github.com/jinzhu/gorm"
	// Register postgres instance
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

// LoadDB loads db initails
func LoadDB() {
	var err error
	conn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		config.Cfg.DBHOST, config.Cfg.DBPORT, config.Cfg.DBUSER, config.Cfg.DBNAME, config.Cfg.DBPASS)
	log.Println("conn---", conn)
	db, err = gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal("Failed to setup db connection, ", err)
	}
	if err = db.DB().Ping(); err != nil {
		log.Fatal("Failed to ping db connection, ", err)
	}
	db.AutoMigrate(
		&Preset{},
		&ProductDiscovery{},
		&MarketArgument{},
	)
}
