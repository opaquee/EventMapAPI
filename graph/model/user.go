package model

type User struct {
	UUIDKey
	FirstName          string   `json:"firstName"`
	LastName           string   `json:"lastName"`
	Email              string   `json:"email"`
	Username           string   `json:"username"`
	Password           string   `json:"password"`
	ProfilePicturePath string   `json:"profilePicturePath"`
	AttendingEvents    []*Event `json:"attendingEvents" gorm:"many2many:user_events;"`
	OwnedEvents        []*Event `json:"ownedEvents"`
}
