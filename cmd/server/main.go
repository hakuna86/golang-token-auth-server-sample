package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/config"
	"github.com/hakuna86/golang-token-auth-server-sample/repo"
	"github.com/hakuna86/golang-token-auth-server-sample/repo/model"
	"github.com/hakuna86/golang-token-auth-server-sample/request"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// database
	rp, err := repo.NewRepo()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer rp.DB.Close()

	// router
	api := e.Group("/api")
	{
		// public api
		pub := api.Group("/public")
		{
			pub.POST("/signUp", request.SignUp(rp))

		}

		// token
		auth := api.Group("/oauth")
		{
			signIn := auth.Group("/signIn")
			{
				signIn.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
					c.Set("username", username)
					c.Set("password", password)
					return rp.IsUser(username, password), nil
				}))
				signIn.POST("/getToken", request.SignIn(rp))
			}
			auth.POST("/refreshToken", request.RefreshToken(rp))
		}

		// private api
		restricted := api.Group("/restricted")
		{
			restricted.Use(middleware.JWT(config.JwtTokenString))
			restricted.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					user := c.Get("user").(*jwt.Token)
					claims := user.Claims.(jwt.MapClaims)
					username := claims["username"].(string)
					u := rp.GetUser(&model.User{Email: username})

					fmt.Print("=============================", username, u.Auth)

					if u.Auth.AccessToken == user.Raw {
						return next(c)
					}
					return errors.New("Access Token is not matched")
				}
			})
			restricted.GET("", request.Restricted())
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}
