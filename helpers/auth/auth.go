package auth

import (
	"context"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/opaquee/EventMapAPI/graph/model"
	"github.com/opaquee/EventMapAPI/helpers/jwt"
	"github.com/opaquee/EventMapAPI/helpers/users"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")

			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			tokenStr := header
			username, err := jwt.ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			userFromDB, err := users.GetUserByUsername(username, db)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user := model.User{
				Username:  username,
				UUIDKey:   userFromDB.UUIDKey,
				Email:     userFromDB.Email,
				FirstName: userFromDB.FirstName,
				LastName:  userFromDB.LastName,
			}

			ctx := context.WithValue(r.Context(), userCtxKey, &user)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) *model.User {
	raw, _ := ctx.Value(userCtxKey).(*model.User)
	return raw
}
