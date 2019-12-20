package main

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/hakuna86/golang-token-auth-server-sample/config"
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
	db, err := config.ConnectDatabase()
	if err != nil {
		e.Logger.Fatal(err)
	}

	r := route.NewRoute(db)

	e.Static("/", "static")

	// public query
	pQuery := e.Group("/publicQuery")
	{
		opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.MaxParallelism(20)}
		pQuery.POST("", r.ServeGraphQl(graphql.MustParseSchema(schma.Schema, schma.NewResolver(db), opts...), false))
	}

	query := e.Group("/query")
	{
		query.Use(middleware.JWT(config.JwtTokenString))
		query.Use(r.AuthMiddleWare)
		opts := []graphql.SchemaOpt{graphql.UseFieldResolvers(), graphql.MaxParallelism(20)}
		query.POST("", r.ServeGraphQl(graphql.MustParseSchema(schma.Schema, schma.NewResolver(db), opts...), true))
	}

	// router
	api := e.Group("/api")
	{
		// public api
		pub := api.Group("/public")
		{
			pub.POST("/signUp", r.SignUp())

		}

		// token
		auth := api.Group("/oauth")
		{
			signIn := auth.Group("/signIn")
			{
				signIn.Use(middleware.BasicAuth(r.AuthBasicMiddleWare))
				signIn.POST("/getToken", r.SignIn())
			}
			auth.POST("/refreshToken", r.RefreshToken())
		}

		// private api
		restricted := api.Group("/restricted")
		{
			restricted.Use(middleware.JWT(config.JwtTokenString))
			restricted.Use(r.AuthMiddleWare)
			restricted.GET("", r.Restricted())
			restricted.GET("/signOut", r.SingOut())
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
}
