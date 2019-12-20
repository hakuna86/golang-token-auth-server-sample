package schma

import (
	"context"
	"github.com/hakuna86/golang-token-auth-server-sample/ent"
	"github.com/hakuna86/golang-token-auth-server-sample/model"
)

const Schema = `
	schema {
		query: Query
	}
	
	type Query {
		user(email: String! Password: String!): User!
	}
	
	type User {
		username: String!
		email: String!
		role: String!
	}
`

type User struct {
	Username string
	Email    string
	Role     string
}

type Resolver struct {
	db *ent.Client
}

func NewResolver(client *ent.Client) *Resolver {
	return &Resolver{
		db: client,
	}
}

func (r *Resolver) User(ctx context.Context, args struct{ Email, Password string }) (User, error) {
	u := &model.User{Eamil: args.Email, Password: args.Password}
	gU, err := u.FindUser(r.db)
	if err != nil {
		return User{}, err
	}
	return User{
		Username: gU.Username,
		Email:    gU.Eamil,
		Role:     gU.Role,
	}, nil
}
