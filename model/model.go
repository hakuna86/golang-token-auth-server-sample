package model

import (
	"context"
	"fmt"
	"github.com/hakuna86/golang-token-auth-server-sample/ent"
	"github.com/hakuna86/golang-token-auth-server-sample/ent/auth"
	"github.com/hakuna86/golang-token-auth-server-sample/ent/user"
	"github.com/labstack/echo"
)

const (
	Admin  = "admin"
	Member = "member"
)

type User struct {
	Eamil    string `json:"email" form:"email" query:"email"`
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
	Role     string `json:"role" form:"role" query:"role"`
}

func (u *User) String() string {
	return fmt.Sprintf("Usernanme: %s, Email: %s, Role : %s", u.Username, u.Eamil, u.Role)
}

func (u *User) CreateUser(client *ent.Client, e echo.Context) error {
	if user, err := client.User.
		Create().
		SetEmail(u.Eamil).
		SetUsername(u.Username).
		SetPassword(u.Password).
		SetRole(u.Role).
		Save(context.Background()); err != nil {
		return err
	} else {
		e.Logger().Debug("user was created: ", user)
	}
	return nil
}

func (u *User) FindUser(client *ent.Client) (*User, error) {
	fu, err := client.User.
		Query().
		Where(
			user.Email(u.Eamil),
			user.Password(u.Password),
		).Only(context.Background())
	if err != nil {
		return nil, err
	}

	return &User{
		Eamil:    fu.Email,
		Username: fu.Username,
		Password: fu.Password,
		Role:     fu.Role,
	}, nil
}

func (u *User) FindByEmail2(client *ent.Client) (*User, error) {
	fu, err := client.User.
		Query().
		Where(
			user.Email(u.Eamil),
		).Only(context.Background())

	if err != nil {
		return nil, err
	}
	return &User{
		Eamil:    fu.Email,
		Username: fu.Username,
		Password: fu.Password,
		Role:     fu.Role,
	}, nil
}

func (u *User) FindByEmail(client *ent.Client) (*ent.User, error) {
	return client.User.
		Query().
		Where(
			user.Email(u.Eamil),
		).Only(context.Background())
}

type Auth struct {
	Access  string `json:"access" form:"access" query:"access"`
	Refresh string `json:"refresh" form:"refresh" query:"refresh"`
}

func (u *Auth) CreateOrUpdateAuthToken(user *User, client *ent.Client) error {
	ctx := context.Background()
	fu, err := user.FindByEmail(client)
	if err != nil {
		return err
	}

	isEx, err := fu.QueryAuth().Exist(ctx)
	if err != nil {
		return err
	}

	if isEx {
		_, err = client.Auth.
			Update().
			SetAccessToken(u.Access).
			SetRefreshToken(u.Refresh).
			Save(ctx)
		return err
	} else {
		_, err = client.Auth.
			Create().
			SetUser(fu).
			SetAccessToken(u.Access).
			SetRefreshToken(u.Refresh).
			Save(ctx)
		return err
	}
}

func (u *Auth) DeleteAuth(client *ent.Client) error {
	_, err := client.Auth.
		Delete().
		Where(
			auth.AccessToken(u.Access),
		).Exec(context.Background())
	return err
}

func (u *Auth) FindUserByAccesstoken(client *ent.Client) (*User, error) {
	ctx := context.Background()
	a, err := client.Auth.
		Query().
		Where(
			auth.AccessToken(u.Access),
		).Only(ctx)
	if err != nil {
		return nil, err
	}
	eu, err := a.QueryUser().Only(ctx)
	if err != nil {
		return nil, err
	}
	return &User{
		Eamil:    eu.Email,
		Username: eu.Username,
		Password: eu.Password,
		Role:     eu.Role,
	}, nil
}
