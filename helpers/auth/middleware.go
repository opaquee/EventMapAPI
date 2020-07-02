package auth

import (
	"context"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/opaquee/EventMapAPI/graph/model"
	"github.com/opaquee/EventMapAPI/helpers/userhelper"
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
			username, err := ParseToken(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			user := model.User{Username: username}
			id, err := userhelper.GetIdByUsername(username)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			user.UUIDKey.ID = id
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
