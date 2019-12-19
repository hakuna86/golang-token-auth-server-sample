package main

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/hakuna86/golang-token-auth-server-sample/config"
	"github.com/hakuna86/golang-token-auth-server-sample/request"
	"github.com/hakuna86/golang-token-auth-server-sample/request/gq"
	"github.com/hakuna86/golang-token-auth-server-sample/request/model"
	mw "github.com/hakuna86/golang-token-auth-server-sample/request/nw"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// database
	dbClient, err := config.ConnectDatabase()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Static("/", "static")

	query := e.Group("/query")
	{
		query.Use(middleware.JWT(config.JwtTokenString))
		query.Use(mw.AuthMiddleWare(dbClient))
		opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.MaxParallelism(20)}
		h := &relay.Handler{Schema: graphql.MustParseSchema(gq.Schema, &gq.Resolver{}, opts...)}
		query.POST("", echo.WrapHandler(h))
	}

	// router
	api := e.Group("/api")
	{
		// public api
		pub := api.Group("/public")
		{
			pub.POST("/signUp", request.SignUp(dbClient))

		}

		// token
		auth := api.Group("/oauth")
		{
			signIn := auth.Group("/signIn")
			{
				signIn.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
					c.Set("username", username)
					c.Set("password", password)
					u := &model.User{Eamil: username, Password: password}
					if _, err := u.FindUser(dbClient); err != nil {
						return false, err
					}
					return true, nil
				}))
				signIn.POST("/getToken", request.SignIn(dbClient))
			}
			auth.POST("/refreshToken", request.RefreshToken(dbClient))
		}

		// private api
		restricted := api.Group("/restricted")
		{
			restricted.Use(middleware.JWT(config.JwtTokenString))
			restricted.Use(mw.AuthMiddleWare(dbClient))
			restricted.GET("", request.Restricted())
			restricted.GET("/signOut", request.SingOut(dbClient))
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}
