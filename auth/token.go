package auth

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"github.com/sisukas/commons/http_utils"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type JWTConfig struct {
	Secret        string
	ExpirySeconds int32
	Aud           string
}
type JWTUtil struct {
	config *JWTConfig
}

func NewJWTUtil(conf *JWTConfig) *JWTUtil {
	return &JWTUtil{config: conf}
}

func (jwu *JWTUtil) GetUserIDFromJWT(r *http.Request) (string, error) {

	if os.Getenv("SISUKAS_ENVIRONMENT") == "development" {

		userID := os.Getenv("SISUKAS_LOGGED_IN_USER_FOR_TESTING")

		if len(userID) > 6 {
			return userID, nil
		}
	}

	strtoken, err := extractBearerToken(r)
	if err != nil || len(strtoken) < 6 {
		log.Printf("Bearer token is empty. Checking cookies ...")

		cookie, err := r.Cookie("simtok")
		if err != nil {
			log.Printf("Error accessing cookie %v", err)
			return "", http_utils.UnauthorizedError("This endpoint requires a valid Bearer token")
		}
		log.Printf("Extracted token from cookie %s ", strtoken)
		strtoken = cookie.Value
	}

	token, err := jwu.ParseJWTClaims(strtoken)
	if err != nil {
		return "", err
	}

	return token.Subject, nil
}

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", http_utils.UnauthorizedError("This endpoint requires a Bearer token")
	}

	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", http_utils.UnauthorizedError("This endpoint requires a Bearer token")
	}

	return matches[1], nil
}

func (jwu *JWTUtil) ParseJWTClaims(bearer string) (*jwt.StandardClaims, error) {
	secret := jwu.config.Secret
	myclaims := &jwt.StandardClaims{}
	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	_, err := p.ParseWithClaims(bearer, myclaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Parsing JWT claims")
	}

	return myclaims, nil
}

func (jwu *JWTUtil) getJWTExpiryDuration() int {

	duration := int(jwu.config.ExpirySeconds)
	if duration < 100 {
		return 3600
	}
	return int(duration)
}

func (jwu *JWTUtil) GenerateAccessToken(user_id string) (string, error) {

	secret := jwu.config.Secret

	expiresIn := time.Second * time.Duration(jwu.getJWTExpiryDuration())

	claims := &jwt.StandardClaims{
		Subject:   user_id,
		Audience:  jwu.config.Aud,
		ExpiresAt: time.Now().Add(expiresIn).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
