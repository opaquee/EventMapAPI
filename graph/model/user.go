package model

type User struct {
	UUIDKey
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	Email              string `json:"email"`
	Username           string `json:"username"`
	Password           string `json:"password"`
	ProfilePicturePath string `json:"profilePicturePath"`
}
