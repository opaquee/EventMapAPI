package dbconn

import "github.com/jinzhu/gorm"

func Open() (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", "host=db port=5432 dbname=postgres user=user password=secret sslmode=disable")
	return db, err
}
