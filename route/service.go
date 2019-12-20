package route

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/config"
	"github.com/hakuna86/golang-token-auth-server-sample/model"
	"time"
)

var (
	authTokenExpir    = time.Now().Add(time.Minute * 10).Unix()
	refreshTokenExpir = time.Now().Add(time.Hour * 24).Unix()
)

func getToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return config.JwtTokenString, nil
	})
}

func makeToken(user *model.User) (*model.Auth, error) {
	// Create Auth token
	token := jwt.New(jwt.SigningMethodHS256)
	//// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["email"] = user.Eamil
	claims["role"] = user.Role
	claims["exp"] = authTokenExpir

	// Generate encoded token and send it as response.
	t, err := token.SignedString(config.JwtTokenString)
	if err != nil {
		return nil, err
	}

	// Create Auth token
	retoken := jwt.New(jwt.SigningMethodHS256)
	//// Set claims
	reClaims := retoken.Claims.(jwt.MapClaims)
	reClaims["refresh"] = true
	reClaims["exp"] = refreshTokenExpir

	// Generate encoded token and send it as response.
	rt, err := retoken.SignedString(config.JwtTokenString)
	if err != nil {
		return nil, err
	}
	return &model.Auth{t, rt}, nil
}
