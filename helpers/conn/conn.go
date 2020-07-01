package conn

import "github.com/jinzhu/gorm"

var DB *gorm.DB

func OpenDB() error {
	db, err := gorm.Open("postgres", "host=db port=5432 dbname=postgres user=user password=secret sslmode=disable")
	DB = db //Save DB handle
	return err
}
