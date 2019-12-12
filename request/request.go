package request

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/repo"
	"github.com/hakuna86/golang-token-auth-server-sample/repo/model"
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
func SignUp(r *repo.Repo) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(model.User)
		if err := c.Bind(u); err != nil {
			return err
		}
		u.Role = "ADMIN"
		if err := r.CreateUser(u); err != nil {
			return c.JSON(http.StatusOK, FailResponse(err))
		}
		return c.JSON(http.StatusOK, SuccesResponse("", nil))
	}
}

func SignIn(r *repo.Repo) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, ok := c.Get("username").(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectInfoError))
		}
		user := r.GetUser(&model.User{Email: username})
		token, err := makeToken(user)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		user.Auth.Email = user.Email
		user.Auth.RefreshToken = token.Refresh
		user.Auth.AccessToken = token.Access

		err = r.UpdateUser(user)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}
		return c.JSON(http.StatusOK, SuccesResponse("Successful SignIn", map[string]string{
			"access_token":  token.Access,
			"refresh_token": token.Refresh,
		}))
	}
}

func SingOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}

func Restricted() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		return c.JSON(http.StatusOK, SuccesResponse("token info", claims))
	}
}

func RefreshToken(r *repo.Repo) echo.HandlerFunc {
	type refreshReq struct {
		RefreshToken string `json:"refreshToken"`
		AuthToken    string `json:"authToken"`
	}

	return func(c echo.Context) error {
		req := new(refreshReq)
		if err := c.Bind(&req); err != nil {
			return err
		}

		_, err := getToken(req.RefreshToken)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		authToken, _ := getToken(req.AuthToken)
		authClaims := authToken.Claims.(jwt.MapClaims)
		username, ok := authClaims["username"].(string)
		if !ok {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectAccessToken))
		}

		user := r.GetUser(&model.User{Email: username})
		if user != nil && user.Auth.RefreshToken == req.RefreshToken {
			token, err := makeToken(user)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, FailResponse(err))
			}

			user.Auth.Email = user.Email
			user.Auth.RefreshToken = token.Refresh
			user.Auth.AccessToken = token.Access

			err = r.UpdateUser(user)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, FailResponse(err))
			}

			return c.JSON(http.StatusOK, SuccesResponse("Successful Refresh Token", map[string]string{
				"access_token":  token.Access,
				"refresh_token": token.Refresh,
			}))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectAccessToken))
		}
	}
}
