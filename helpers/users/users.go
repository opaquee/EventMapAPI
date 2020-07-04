package users

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/opaquee/EventMapAPI/graph/model"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetIdByUsername(username string, db *gorm.DB) (id uuid.UUID, err error) {
	user := model.User{
		Username: username,
	}
	if err := db.Where(&user).First(&user).Error; err != nil {
		log.Print("Failed to get user by username")
		return uuid.UUID{}, err
	}
	return user.UUIDKey.ID, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
