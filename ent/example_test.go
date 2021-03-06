// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"log"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// dsn for the database. In order to run the tests locally, run the following command:
//
//	 ENT_INTEGRATION_ENDPOINT="root:pass@tcp(localhost:3306)/test?parseTime=True" go test -v
//
var dsn string

func ExampleArticle() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the article's edges.

	// create article vertex with its edges.
	a := client.Article.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetTitle("string").
		SetBody("string").
		SetStar(1).
		SaveX(ctx)
	log.Println("article created:", a)

	// query edges.

	// Output:
}
func ExampleAuth() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the auth's edges.

	// create auth vertex with its edges.
	a := client.Auth.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetAccessToken("string").
		SetRefreshToken("string").
		SaveX(ctx)
	log.Println("auth created:", a)

	// query edges.

	// Output:
}
func ExampleUser() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the user's edges.
	a0 := client.Auth.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetAccessToken("string").
		SetRefreshToken("string").
		SaveX(ctx)
	log.Println("auth created:", a0)
	a1 := client.Article.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetTitle("string").
		SetBody("string").
		SetStar(1).
		SaveX(ctx)
	log.Println("article created:", a1)

	// create user vertex with its edges.
	u := client.User.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetEmail("string").
		SetUsername("string").
		SetPassword("string").
		SetRole("string").
		SetAuth(a0).
		AddArticles(a1).
		SaveX(ctx)
	log.Println("user created:", u)

	// query edges.
	a0, err = u.QueryAuth().First(ctx)
	if err != nil {
		log.Fatalf("failed querying auth: %v", err)
	}
	log.Println("auth found:", a0)

	a1, err = u.QueryArticles().First(ctx)
	if err != nil {
		log.Fatalf("failed querying articles: %v", err)
	}
	log.Println("articles found:", a1)

	// Output:
}
