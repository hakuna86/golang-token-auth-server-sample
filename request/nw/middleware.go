// middleware

package mw

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/ent"
	"github.com/hakuna86/golang-token-auth-server-sample/request"
	"github.com/hakuna86/golang-token-auth-server-sample/request/model"
	"github.com/labstack/echo"
	"net/http"
)

func AuthMiddleWare(client *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			access := c.Get("user").(*jwt.Token)
			a := &model.Auth{Access: access.Raw}
			user, err := a.FindUserByAccesstoken(client)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err)
			}
			claims := access.Claims.(jwt.MapClaims)
			email := claims["email"].(string)
			if email == user.Eamil {
				return next(c)
			}
			return c.JSON(http.StatusUnauthorized, request.IncorrectAccessToken)
		}
	}
}
