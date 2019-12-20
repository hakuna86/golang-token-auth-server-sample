package route

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/ent"
	"github.com/hakuna86/golang-token-auth-server-sample/model"
	"github.com/labstack/echo"
	"net/http"
)

var (
	IncorrectInfoError   = errors.New("IncorrectInfoError")
	IncorrectAccessToken = errors.New("IncorrectAccessToken")
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccesResponse(message string, data interface{}) *Response {
	return &Response{true, message, data}
}

func FailResponse(err error) *Response {
	return &Response{false, err.Error(), nil}
}

// Member register
func SignUp(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(model.User)
		if err := c.Bind(u); err != nil {
			return err
		}
		u.Role = model.Admin
		if err := u.CreateUser(client, c); err != nil {
			return c.JSON(http.StatusConflict, FailResponse(err))
		}
		return c.JSON(http.StatusCreated, SuccesResponse(u.String(), nil))
	}
}

func SignIn(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectInfoError))
		}

		password, ok := c.Get("password").(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectInfoError))
		}

		u := &model.User{Eamil: username, Password: password}
		user, err := u.FindUser(client)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		auth, err := makeToken(user)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		if err := auth.CreateOrUpdateAuthToken(user, client); err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		return c.JSON(http.StatusOK, SuccesResponse("Successful SignIn", auth))
	}
}

func SingOut(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		a := &model.Auth{Access: token.Raw}
		if err := a.DeleteAuth(client); err != nil {
			return c.JSON(http.StatusOK, FailResponse(err))
		}
		return c.JSON(http.StatusOK, SuccesResponse("SignOut Successed", nil))
	}
}

func Restricted() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		return c.JSON(http.StatusOK, SuccesResponse("token info", claims))
	}
}

func RefreshToken(client *ent.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		au := new(model.Auth)
		if err := c.Bind(au); err != nil {
			return err
		}
		_, err := getToken(au.Refresh)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}
		authToken, _ := getToken(au.Access)
		authClaims := authToken.Claims.(jwt.MapClaims)
		email, ok := authClaims["email"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectAccessToken))
		}
		user, err := au.FindUserByAccesstoken(client)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectAccessToken))
		}
		if user.Eamil == email {
			auth, err := makeToken(user)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, FailResponse(err))
			}
			if err := auth.CreateOrUpdateAuthToken(user, client); err != nil {
				return c.JSON(http.StatusUnauthorized, FailResponse(err))
			}
			return c.JSON(http.StatusOK, SuccesResponse("Successful token refresh", auth))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectAccessToken))
		}
	}
}
