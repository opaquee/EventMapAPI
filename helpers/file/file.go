package file

import (
	"errors"
	"strings"

	"github.com/opaquee/EventMapAPI/graph/model"
)

func ValidImageFile(filename string) (bool, error) {
	var validImgExt = []string{"bmp", "jpeg", "jpg", "png"}

	splitFile := strings.Split(filename, ".")
	if len(splitFile) != 2 {
		return false, errors.New("File name is incorrect")
	}
	fileExt := strings.ToLower(splitFile[1])

	valid := false
	for _, ext := range validImgExt {
		if fileExt == strings.ToLower(ext) {
			valid = true
			break
		}
	}

	return valid, nil
}

func NewFileName(filename string, user *model.User) string {
	return user.UUIDKey.ID.String() + "." + strings.Split(filename, ".")[1]
}
