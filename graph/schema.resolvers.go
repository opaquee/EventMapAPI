package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opaquee/EventMapAPI/graph/generated"
	"github.com/opaquee/EventMapAPI/graph/model"
	"github.com/opaquee/EventMapAPI/helpers/jwt"
	"github.com/opaquee/EventMapAPI/helpers/users"
)

func (r *eventResolver) ID(ctx context.Context, obj *model.Event) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	hashedPassword, err := users.HashPassword(input.Password)
	if err != nil {
		log.Print("error hashing password")
		return "", err
	}

	user := model.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Username:  input.Username,
		Password:  hashedPassword,
	}

	if err := r.DB.Create(&user).Error; err != nil {
		log.Print("failed to create user")
		return "", err
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, userID int) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	user := model.User{
		Username: input.Username,
		Password: input.Password,
	}

	if users.Authenticate(&user, r.DB) == false {
		return "", errors.New("wrong username or password")
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
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

func (r *userResolver) ID(ctx context.Context, obj *model.User) (string, error) {
	return obj.UUIDKey.ID.String(), nil
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type eventResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
