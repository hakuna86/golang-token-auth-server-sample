package route

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/graph-gophers/graphql-go"
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

type Route struct {
	db *ent.Client
}

func NewRoute(client *ent.Client) *Route {
	return &Route{db: client}
}

func (r *Route) ServeGraphQl(schema *graphql.Schema, auth bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		ctx := c.Request().Context()
		if auth {
			ctx = context.WithValue(c.Request().Context(), "token", c.Get("user").(*jwt.Token))
		}
		response := schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.Blob(http.StatusOK, "application/json", responseJSON)
	}
}

// Member register
func (r *Route) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		u := new(model.User)
		if err := c.Bind(u); err != nil {
			return err
		}
		u.Role = model.Admin
		if err := u.CreateUser(r.db, c); err != nil {
			return c.JSON(http.StatusConflict, FailResponse(err))
		}
		return c.JSON(http.StatusCreated, SuccesResponse(u.String(), nil))
	}
}

func (r *Route) SignIn() echo.HandlerFunc {
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
		user, err := u.FindUser(r.db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		auth, err := makeToken(user)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		if err := auth.CreateOrUpdateAuthToken(user, r.db); err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(err))
		}

		return c.JSON(http.StatusOK, SuccesResponse("Successful SignIn", auth))
	}
}

func (r *Route) SingOut() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user").(*jwt.Token)
		a := &model.Auth{Access: token.Raw}
		if err := a.DeleteAuth(r.db); err != nil {
			return c.JSON(http.StatusOK, FailResponse(err))
		}
		return c.JSON(http.StatusOK, SuccesResponse("SignOut Successed", nil))
	}
}

func (r *Route) Restricted() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		return c.JSON(http.StatusOK, SuccesResponse("token info", claims))
	}
}

func (r *Route) RefreshToken() echo.HandlerFunc {
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
		user, err := au.FindUserByAccesstoken(r.db)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectAccessToken))
		}
		if user.Eamil == email {
			auth, err := makeToken(user)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, FailResponse(err))
			}
			if err := auth.CreateOrUpdateAuthToken(user, r.db); err != nil {
				return c.JSON(http.StatusUnauthorized, FailResponse(err))
			}
			return c.JSON(http.StatusOK, SuccesResponse("Successful token refresh", auth))
		} else {
			return c.JSON(http.StatusUnauthorized, FailResponse(IncorrectAccessToken))
		}
	}
}
