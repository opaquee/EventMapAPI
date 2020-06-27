package model

import "time"

type Event struct {
	UUIDKey
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	AddressLine1 string    `json:"addressLine1"`
	AddressLine2 string    `json:"addressLine2"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Zip          int       `json:"zip"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	Users        []*User   `json:"users"`
}
