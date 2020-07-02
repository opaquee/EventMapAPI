package conn

import "github.com/jinzhu/gorm"

func OpenDB() (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", "host=db port=5432 dbname=postgres user=user password=secret sslmode=disable")
	return db, err
}
