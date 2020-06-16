package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opaquee/EventMapAPI/graph/generated"
	"github.com/opaquee/EventMapAPI/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, userID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateEvent(ctx context.Context, input model.NewEvent) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteEvent(ctx context.Context, eventID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddUserProfilePicture(ctx context.Context, profilePicture graphql.Upload) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveUserProfilePicture(ctx context.Context, userID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddUserToEvent(ctx context.Context, userID int, eventID int) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveUserFromEvent(ctx context.Context, userID int, eventID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAllNearbyEvents(ctx context.Context, latitude float64, longitude float64) ([]*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetEventUsers(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetEventByID(ctx context.Context, eventID int) (*model.Event, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUserByID(ctx context.Context, userID int) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }


