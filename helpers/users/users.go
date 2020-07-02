package users

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/opaquee/EventMapAPI/graph/model"
	"golang.org/x/crypto/bcrypt"
)

func GetIdByUsername(username string, db *gorm.DB) (id string, err error) {
	user := model.User{
		Username: username,
	}
	if err := db.Where(&user).First(&user).Error; err != nil {
		log.Print("Failed to get user by username")
		return "", err
	}
	return user.UUIDKey.ID.String(), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
