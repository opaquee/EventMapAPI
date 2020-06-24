package db_boilerplate

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func migrateDB(db *DB) {
	db.AutoMigrate(&User{}, &Event{})
}
