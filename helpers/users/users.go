package users

import (
	"errors"

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
		return uuid.UUID{}, err
	}

	return user.UUIDKey.ID, nil
}

func GetUserByUsername(username string, db *gorm.DB) (user *model.User, err error) {
	user = &model.User{
		Username: username,
	}

	if err := db.Where(&user).First(&user).Error; err != nil {
		return nil, err
	}

	//Omit password, pfp information, etc. Meant for use in auth middleware
	userFields := &model.User{
		UUIDKey:   user.UUIDKey,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return userFields, nil
}

func Authenticate(incomingUser *model.User, db *gorm.DB) (correct bool, err error) {
	userFromDB := model.User{}

	if err = db.Where(model.User{
		Username: incomingUser.Username,
	}).First(&userFromDB).Error; err != nil {
		return false, err
	}

	return CheckPasswordHash(incomingUser.Password, userFromDB.Password), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Duplicate(incomingUser *model.User, db *gorm.DB) (err error) {
	var user model.User

	err = db.Where(model.User{
		Username: incomingUser.Username,
	}).Or(model.User{
		Email: incomingUser.Email,
	}).First(&user).Error

	if err != nil && gorm.IsRecordNotFoundError(err) == false {
		return err
	}

	if gorm.IsRecordNotFoundError(err) == false {
		return errors.New("username or email is duplicate")
	}

	return nil
}
