package request

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/config"
	"github.com/hakuna86/golang-token-auth-server-sample/repo"
	"github.com/hakuna86/golang-token-auth-server-sample/repo/model"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type Response struct {
	Status bool `json:"status"`
	Message string `json:"message"`
}

func SuccesResponse(message string) *Response {
	return &Response{true, message}
}

func FailResponse(err error) *Response {
	return &Response{false, err.Error()}
}

// Member register
func SignUp(r *repo.Repo) echo.HandlerFunc{
	return func(c echo.Context) error {
		u := new(model.User)
		if err := c.Bind(u); err != nil {
			return err
		}
		u.Role = "ADMIN"

		if err := r.CreateUser(u); err != nil {
			return c.JSON(http.StatusOK, FailResponse(err))
		}
		return c.JSON(http.StatusOK, SuccesResponse(""))
	}
}

func SignIn(r *repo.Repo) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.Get("username").(string)
		user := r.GetUser(username)
		fmt.Println("===========", user.Role)

		// Create token
		token := jwt.New(jwt.SigningMethodHS256)
		//// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = "username"
		claims["role"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString(config.JwtSignString)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}

func SingOut()  echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func Restricted() echo.HandlerFunc {
	return func (c echo.Context) error {
		//user := c.Get("user").(*jwt.Token)
		//claims := user.Claims.(jwt.MapClaims)
		//name := claims["name"].(string)
		return c.String(http.StatusOK, "Welcome !")
	}
}
