package schma

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/hakuna86/golang-token-auth-server-sample/ent"
	"github.com/hakuna86/golang-token-auth-server-sample/model"
)

var (
	tokenIsNil = errors.New("Token is Nil")
)

const Schema = `
	schema {
		query: Query
	}
	
	type Query {
		user(): User!
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

func (r *Resolver) User(ctx context.Context) (User, error) {
	token, ok := ctx.Value("token").(*jwt.Token)
	if !ok {
		return User{}, tokenIsNil
	}
	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)
	u := &model.User{Eamil: email}
	gU, err := u.FindByEmail2(r.db)
	if err != nil {
		return User{}, err
	}
	return User{
		Username: gU.Username,
		Email:    gU.Eamil,
		Role:     gU.Role,
	}, nil
}
