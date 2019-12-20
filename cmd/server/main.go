package main

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/hakuna86/golang-token-auth-server-sample/config"
	"github.com/hakuna86/golang-token-auth-server-sample/model"
	"github.com/hakuna86/golang-token-auth-server-sample/model/schma"
	"github.com/hakuna86/golang-token-auth-server-sample/route"
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

	// public query
	pQuery := e.Group("/publicQuery")
	{
		opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.MaxParallelism(20)}
		h := &relay.Handler{Schema: graphql.MustParseSchema(schma.Schema, schma.NewResolver(dbClient), opts...)}
		pQuery.POST("", echo.WrapHandler(h))
	}

	query := e.Group("/query")
	{
		query.Use(middleware.JWT(config.JwtTokenString))
		query.Use(route.AuthMiddleWare(dbClient))
		opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.MaxParallelism(20)}
		h := &relay.Handler{Schema: graphql.MustParseSchema(schma.Schema, schma.NewResolver(dbClient), opts...)}
		query.POST("", echo.WrapHandler(h))
	}

	// router
	api := e.Group("/api")
	{
		// public api
		pub := api.Group("/public")
		{
			pub.POST("/signUp", route.SignUp(dbClient))

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
				signIn.POST("/getToken", route.SignIn(dbClient))
			}
			auth.POST("/refreshToken", route.RefreshToken(dbClient))
		}

		// private api
		restricted := api.Group("/restricted")
		{
			restricted.Use(middleware.JWT(config.JwtTokenString))
			restricted.Use(route.AuthMiddleWare(dbClient))
			restricted.GET("", route.Restricted())
			restricted.GET("/signOut", route.SingOut(dbClient))
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}
