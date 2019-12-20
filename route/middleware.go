// middleware

package route

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/model"
	"github.com/labstack/echo"
	"net/http"
)

func (r *Route) AuthMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		access := c.Get("user").(*jwt.Token)
		a := &model.Auth{Access: access.Raw}
		user, err := a.FindUserByAccesstoken(r.db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}
		claims := access.Claims.(jwt.MapClaims)
		email := claims["email"].(string)
		if email == user.Eamil {
			return next(c)
		}
		return c.JSON(http.StatusUnauthorized, IncorrectAccessToken)
	}
}

func (r *Route) AuthBasicMiddleWare(username, password string, c echo.Context) (bool, error) {
	c.Set("username", username)
	c.Set("password", password)
	u := &model.User{Eamil: username, Password: password}
	if _, err := u.FindUser(r.db); err != nil {
		return false, err
	}
	return true, nil
}
