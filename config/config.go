package config

import (
	"context"
	"fmt"
	"github.com/hakuna86/golang-token-auth-server-sample/ent"
	_ "github.com/lib/pq"
)

//_ "github.com/mattn/go-sqlite3" --> using sqlite
//_ "github.com/lib/pq" --> usin postgres

var (
	JwtTokenString   = []byte("secret")
	connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "localhost", "5432", "postgres", "11111", "test")
)

func ConnectDatabase() (*ent.Client, error) {
	client, err := ent.Open("postgres", connectionString)
	//client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, err
	}
	// run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}
	return client, nil
}
