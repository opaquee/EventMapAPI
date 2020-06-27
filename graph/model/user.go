package model

import "github.com/99designs/gqlgen/graphql"

type User struct {
	UUIDKey
	FirstName      string          `json:"firstName"`
	LastName       string          `json:"lastName"`
	Email          string          `json:"email"`
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	ProfilePicture *graphql.Upload `json:"profilePicture"`
	Events         []*Event        `json:"events"`
}
