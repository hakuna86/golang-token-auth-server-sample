package main

import (
	"github.com/hakuna86/golang-token-auth-server-sample/config"
	"github.com/hakuna86/golang-token-auth-server-sample/repo"
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
	r, err := repo.NewRepo()
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer r.DB.Close()

	// router
	api := e.Group("/api")
	{
		// public api
		pub := api.Group("/public")
		{
			pub.POST("/signUp", request.SignUp(r))
		}

		// token
		auth := api.Group("/oauth")
		{
			auth.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
				c.Set("username", username)
				c.Set("password", password)
				return r.IsUser(username, password), nil
			}))
			auth.POST("/getToken", request.SignIn(r))
		}

		// private api
		r := api.Group("/restricted")
		{
			r.Use(middleware.JWT(config.JwtSignString))
			r.GET("", request.Restricted())
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}