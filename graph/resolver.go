package graph

import (
	"sync"

	"github.com/jinzhu/gorm"
	"github.com/opaquee/EventMapAPI/graph/model"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	MU sync.Mutex
	//Observers map[int](map[string]chan *model.Event)
	Observers map[string]chan *model.Event
	DB        *gorm.DB
}
