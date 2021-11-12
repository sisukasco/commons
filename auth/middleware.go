package auth

import (
	"context"
	"log"
	"net/http"
	"github.com/sisukas/commons/http_utils"
)

//UserIDKeyType is the key type for setting the context and to avoid collisions
type UserIDKeyType string

const (
	//UserIDKey is the key used to set the context value for User ID
	UserIDKey UserIDKeyType = "SimUserID"
)

type AuthMiddleware struct {
	jwtx *JWTUtil
}

func NewAuthMiddleware(jwtConf *JWTConfig) *AuthMiddleware {
	jwtx := NewJWTUtil(jwtConf)
	return &AuthMiddleware{jwtx: jwtx}
}

func (am *AuthMiddleware) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, err := am.jwtx.GetUserIDFromJWT(r)
		if err != nil {
			log.Printf("No UserID in JWT. Error %v", err)
			next.ServeHTTP(w, r)
			return
		}

		log.Printf("Got UserID from JWT Token %s", userID)

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID, err := am.jwtx.GetUserIDFromJWT(r)
		if err != nil {
			err := http_utils.UnauthorizedError("This endpoint requires a valid Bearer token")
			http_utils.SendErrorResponse(err, w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
