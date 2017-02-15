package sqlite

import (
	"github.com/kassisol/twic/storage"
	"github.com/kassisol/twic/storage/driver"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	storage.RegisterDriver("sqlite", New)
}

type Config struct {
	DB *gorm.DB
}

func New(dbFilePath string) (driver.Storager, error) {
	debug := false

	db, err := gorm.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	db.LogMode(debug)

	db.AutoMigrate(&Cert{})
	db.AutoMigrate(&Profile{})

	return &Config{DB: db}, nil
}

func (c *Config) End() {
	c.DB.Close()
}
