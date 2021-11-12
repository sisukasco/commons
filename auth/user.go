package auth

import (
	"github.com/pkg/errors"
	"net/http"
)

func GetUserIDFromContext(r *http.Request) (string, error) {
	ctx := r.Context()
	userv := ctx.Value(UserIDKey)
	if userv == nil {
		return "", errors.New("Context does not have user")
	}
	userID := userv.(string)

	return userID, nil
}
