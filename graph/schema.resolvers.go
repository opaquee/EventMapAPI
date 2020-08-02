package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/opaquee/EventMapAPI/graph/generated"
	"github.com/opaquee/EventMapAPI/graph/model"
	"github.com/opaquee/EventMapAPI/helpers/auth"
	"github.com/opaquee/EventMapAPI/helpers/file"
	"github.com/opaquee/EventMapAPI/helpers/geocode"
	"github.com/opaquee/EventMapAPI/helpers/jwt"
	"github.com/opaquee/EventMapAPI/helpers/users"
	uuid "github.com/satori/go.uuid"
)

func (r *eventResolver) ID(ctx context.Context, obj *model.Event) (string, error) {
	return obj.UUIDKey.ID.String(), nil
}

func (r *eventResolver) Users(ctx context.Context, obj *model.Event) ([]*model.User, error) {
	var users []*model.User

	if err := r.DB.Model(model.Event{
		UUIDKey: obj.UUIDKey,
	}).Association("Users").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *eventResolver) Owner(ctx context.Context, obj *model.Event) (*model.User, error) {
	ownerKey := model.UUIDKey{
		ID: obj.OwnerID,
	}
	owner := &model.User{
		UUIDKey: ownerKey,
	}

	if err := r.DB.Where(owner).First(owner).Error; err != nil {
		return nil, err
	}

	return owner, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	if err := users.Duplicate(&model.User{
		Email:    input.Email,
		Username: input.Username,
	}, r.DB); err != nil {
		return "", err
	}

	hashedPassword, err := users.HashPassword(input.Password)
	if err != nil {
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
		return "", err
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, username string, input model.UpdateUserInput) (*model.User, error) {
	userFromCtx := auth.ForContext(ctx)
	userFromDB, err := users.GetUserByUsername(username, r.DB)
	if err != nil {
		return nil, err
	}

	if err := users.CheckAccess(userFromCtx, userFromDB); err != nil {
		return nil, err
	}

	//if email is duplicate, reject
	if userFromCtx.Email != input.Email {
		if err := users.Duplicate(&model.User{
			Email: input.Email,
		}, r.DB); err != nil {
			return nil, err
		}
	}

	userFromDB.FirstName = input.FirstName
	userFromDB.LastName = input.LastName
	userFromDB.Email = input.Email

	if err := r.DB.Save(userFromDB).Error; err != nil {
		return nil, err
	}

	return userFromDB, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, username string) (bool, error) {
	userFromCtx := auth.ForContext(ctx)
	userFromDB, err := users.GetUserByUsername(username, r.DB)
	if err != nil {
		return false, err
	}

	if err := users.CheckAccess(userFromCtx, userFromDB); err != nil {
		return false, err
	}

	if err := r.DB.Model(model.User{
		UUIDKey: userFromDB.UUIDKey,
	}).Association("AttendingEvents").Clear().Error; err != nil {
		return false, err
	}

	if err := r.DB.Unscoped().Delete(&userFromDB).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.LoginResponse, error) {
	user := model.User{
		Username: input.Username,
		Password: input.Password,
	}

	correctLogin, err := users.Authenticate(&user, r.DB)
	if err != nil {
		return nil, err
	}
	if correctLogin == false {
		return nil, errors.New("incorrect username or password")
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return nil, err
	}

	userFromDB, err := users.GetUserByUsername(user.Username, r.DB)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		User:  userFromDB,
	}, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", errors.New("access denied")
	}

	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*model.Event, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return nil, errors.New("no user information from context. You probably didn't provide a token")
	}

	event := model.Event{
		Name:         input.Name,
		Description:  input.Description,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		State:        input.State,
		Zip:          input.Zip,
		OwnerID:      userFromCtx.ID,
	}

	//Get latitude and longitude from the geocoding api
	if err := geocode.GetLatLng(&event); err != nil {
		return nil, err
	}

	if err := r.DB.Create(&event).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Model(&model.User{
		UUIDKey: userFromCtx.UUIDKey,
	}).Association("OwnedEvents").Append(&event).Error; err != nil {
		return nil, err
	}

	r.MU.Lock()
	for _, observer := range r.Observers[input.Zip] {
		observer <- &event
	}
	r.MU.Unlock()

	return &event, nil
}

func (r *mutationResolver) UpdateEvent(ctx context.Context, eventID string, input model.NewEvent) (*model.Event, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return nil, errors.New("no user information from context. You probably didn't provide a token")
	}

	id, err := uuid.FromString(eventID)
	if err != nil {
		return nil, err
	}
	oldEvent := &model.Event{
		UUIDKey: model.UUIDKey{
			ID: id,
		},
	}

	if err := users.CheckEventOwner(userFromCtx, oldEvent, r.DB); err != nil {
		return nil, err
	}

	newEvent := model.Event{
		Name:         input.Name,
		Description:  input.Description,
		AddressLine1: input.AddressLine1,
		AddressLine2: input.AddressLine2,
		City:         input.City,
		State:        input.State,
		Zip:          input.Zip,
		Latitude:     oldEvent.Latitude,
		Longitude:    oldEvent.Longitude,
		OwnerID:      userFromCtx.UUIDKey.ID,
	}

	id, err = uuid.FromString(eventID)
	if err != nil {
		return nil, err
	}
	newEvent.UUIDKey.ID = id

	//If address is new, get latitude and longitude from the geocoding api
	if oldEvent.AddressLine1 != newEvent.AddressLine1 ||
		oldEvent.City != newEvent.City ||
		oldEvent.State != newEvent.State ||
		oldEvent.Zip != newEvent.Zip {
		if err := geocode.GetLatLng(&newEvent); err != nil {
			return nil, err
		}
	}

	if err := r.DB.Save(&newEvent).Error; err != nil {
		return nil, err
	}

	r.MU.Lock()
	for _, observer := range r.Observers[input.Zip] {
		observer <- &newEvent
	}
	r.MU.Unlock()

	return &newEvent, nil
}

func (r *mutationResolver) DeleteEvent(ctx context.Context, eventID string) (bool, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return false, errors.New("no user information from context. You probably didn't provide a token")
	}

	id, err := uuid.FromString(eventID)
	if err != nil {
		return false, err
	}
	UUIDKey := model.UUIDKey{
		ID: id,
	}
	event := &model.Event{
		UUIDKey: UUIDKey,
	}

	if err := users.CheckEventOwner(userFromCtx, event, r.DB); err != nil {
		return false, err
	}

	if err := r.DB.Model(model.Event{
		UUIDKey: UUIDKey,
	}).Association("Users").Clear().Error; err != nil {
		return false, err
	}

	if err := r.DB.Unscoped().Delete(&model.Event{
		UUIDKey: UUIDKey,
	}).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AddUserProfilePicture(ctx context.Context, profilePicture graphql.Upload) (bool, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return false, errors.New("no user information from context. You probably didn't provide a token")
	}

	if valid, err := file.ValidImageFile(profilePicture.Filename); err != nil || valid == false {
		return false, errors.New("Invalid image file")
	}
	content, err := ioutil.ReadAll(profilePicture.File)
	if err != nil {
		return false, err
	}

	filePath := os.Getenv("APP_VOLUME") + file.NewFileName(profilePicture.Filename, userFromCtx)
	if _, err := os.Create(filePath); err != nil {
		return false, err
	}
	if err := ioutil.WriteFile(filePath, content, 0644); err != nil {
		return false, err
	}

	userFromDB, err := users.GetUserByUsername(userFromCtx.Username, r.DB)
	if err != nil {
		return false, err
	}
	userFromDB.ProfilePicturePath = filePath
	if err := r.DB.Save(userFromDB).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RemoveUserProfilePicture(ctx context.Context) (bool, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return false, errors.New("no user information from context. You probably didn't provide a token")
	}

	splitPath := strings.Split(userFromCtx.ProfilePicturePath, ".")
	if splitPath[0] != os.Getenv("APP_VOLUME")+userFromCtx.UUIDKey.ID.String() {
		return false, errors.New("Access denied, can't delete file at specified path")
	}

	if err := os.Remove(userFromCtx.ProfilePicturePath); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AddUserToEvent(ctx context.Context, eventID string) (bool, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return false, errors.New("no user information from context. You probably didn't provide a token")
	}

	id, err := uuid.FromString(eventID)
	if err != nil {
		return false, err
	}
	r.DB.Model(userFromCtx).Association("AttendingEvents").Append(&model.Event{
		UUIDKey: model.UUIDKey{
			ID: id,
		},
	})

	return true, nil
}

func (r *mutationResolver) RemoveUserFromEvent(ctx context.Context, eventID string) (bool, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return false, errors.New("no user information from context. You probably didn't provide a token")
	}

	id, err := uuid.FromString(eventID)
	if err != nil {
		return false, err
	}
	r.DB.Model(userFromCtx).Association("AttendingEvents").Delete(&model.Event{
		UUIDKey: model.UUIDKey{
			ID: id,
		},
	})

	return true, nil
}

func (r *queryResolver) GetAllNearbyEvents(ctx context.Context, zip int) ([]*model.Event, error) {
	var nearbyEvents []*model.Event

	if err := r.DB.Where(&model.Event{
		Zip: zip,
	}).Find(&nearbyEvents).Error; err != nil {
		return nil, err
	}

	return nearbyEvents, nil
}

func (r *queryResolver) GetEventByID(ctx context.Context, eventID string) (*model.Event, error) {
	id, err := uuid.FromString(eventID)
	if err != nil {
		return nil, err
	}
	eventFromDB := model.Event{
		UUIDKey: model.UUIDKey{
			ID: id,
		},
	}

	if err := r.DB.Where(&eventFromDB).First(&eventFromDB).Error; err != nil {
		return nil, err
	}

	return &eventFromDB, nil
}

func (r *queryResolver) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	id, err := uuid.FromString(userID)
	if err != nil {
		return nil, err
	}
	userFromDB := model.User{
		UUIDKey: model.UUIDKey{
			ID: id,
		},
	}

	if err := r.DB.Where(&userFromDB).First(&userFromDB).Error; err != nil {
		return nil, err
	}

	return &userFromDB, nil
}

func (r *subscriptionResolver) NewEvents(ctx context.Context, zip int, userID string) (<-chan *model.Event, error) {
	observer := make(chan *model.Event, 1)

	//Cleanup empty observer channels
	go func() {
		<-ctx.Done()
		r.MU.Lock()
		delete(r.Observers[zip], userID)
		if len(r.Observers[zip]) == 0 {
			delete(r.Observers, zip)
		}
		r.MU.Unlock()
	}()

	r.MU.Lock()
	if r.Observers[zip] == nil {
		localObservers := make(map[string]chan *model.Event, 1)
		r.Observers[zip] = localObservers
	}
	r.Observers[zip][userID] = observer
	r.MU.Unlock()

	return observer, nil
}

func (r *userResolver) ID(ctx context.Context, obj *model.User) (string, error) {
	return obj.UUIDKey.ID.String(), nil
}

func (r *userResolver) Email(ctx context.Context, obj *model.User) (string, error) {
	return "", errors.New("access denied, email is private to user")
}

func (r *userResolver) Password(ctx context.Context, obj *model.User) (string, error) {
	return "", errors.New("access denied, password is private to user")
}

func (r *userResolver) ProfilePicture(ctx context.Context, obj *model.User) (*model.File, error) {
	userFromCtx := auth.ForContext(ctx)
	if userFromCtx == nil {
		return nil, errors.New("no user information from context. You probably didn't provide a token")
	}

	fileBytes, err := ioutil.ReadFile(userFromCtx.ProfilePicturePath)
	if err != nil {
		return nil, err
	}

	ext := strings.Split(userFromCtx.ProfilePicturePath, ".")[1]

	return &model.File{
		Name:        userFromCtx.Username + "." + ext,
		Content:     string(fileBytes),
		ContentType: "image/" + ext,
	}, nil
}

func (r *userResolver) AttendingEvents(ctx context.Context, obj *model.User) ([]*model.Event, error) {
	var attendingEvents []*model.Event

	if err := r.DB.Model(model.User{
		UUIDKey: obj.UUIDKey,
	}).Association("AttendingEvents").Find(&attendingEvents).Error; err != nil {
		return nil, err
	}

	return attendingEvents, nil
}

func (r *userResolver) OwnedEvents(ctx context.Context, obj *model.User) ([]*model.Event, error) {
	var ownedEvents []*model.Event

	if err := r.DB.Model(model.User{
		UUIDKey: obj.UUIDKey,
	}).Association("OwnedEvents").Find(&ownedEvents).Error; err != nil {
		return nil, err
	}

	return ownedEvents, nil
}

// Event returns generated.EventResolver implementation.
func (r *Resolver) Event() generated.EventResolver { return &eventResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type eventResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
