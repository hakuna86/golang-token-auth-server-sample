// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/hakuna86/golang-token-auth-server-sample/ent/auth"
)

// AuthCreate is the builder for creating a Auth entity.
type AuthCreate struct {
	config
	created_at    *time.Time
	updated_at    *time.Time
	access_token  *string
	refresh_token *string
	user          map[int]struct{}
}

// SetCreatedAt sets the created_at field.
func (ac *AuthCreate) SetCreatedAt(t time.Time) *AuthCreate {
	ac.created_at = &t
	return ac
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (ac *AuthCreate) SetNillableCreatedAt(t *time.Time) *AuthCreate {
	if t != nil {
		ac.SetCreatedAt(*t)
	}
	return ac
}

// SetUpdatedAt sets the updated_at field.
func (ac *AuthCreate) SetUpdatedAt(t time.Time) *AuthCreate {
	ac.updated_at = &t
	return ac
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (ac *AuthCreate) SetNillableUpdatedAt(t *time.Time) *AuthCreate {
	if t != nil {
		ac.SetUpdatedAt(*t)
	}
	return ac
}

// SetAccessToken sets the access_token field.
func (ac *AuthCreate) SetAccessToken(s string) *AuthCreate {
	ac.access_token = &s
	return ac
}

// SetRefreshToken sets the refresh_token field.
func (ac *AuthCreate) SetRefreshToken(s string) *AuthCreate {
	ac.refresh_token = &s
	return ac
}

// SetUserID sets the user edge to User by id.
func (ac *AuthCreate) SetUserID(id int) *AuthCreate {
	if ac.user == nil {
		ac.user = make(map[int]struct{})
	}
	ac.user[id] = struct{}{}
	return ac
}

// SetUser sets the user edge to User.
func (ac *AuthCreate) SetUser(u *User) *AuthCreate {
	return ac.SetUserID(u.ID)
}

// Save creates the Auth in the database.
func (ac *AuthCreate) Save(ctx context.Context) (*Auth, error) {
	if ac.created_at == nil {
		v := auth.DefaultCreatedAt()
		ac.created_at = &v
	}
	if ac.updated_at == nil {
		v := auth.DefaultUpdatedAt()
		ac.updated_at = &v
	}
	if ac.access_token == nil {
		return nil, errors.New("ent: missing required field \"access_token\"")
	}
	if ac.refresh_token == nil {
		return nil, errors.New("ent: missing required field \"refresh_token\"")
	}
	if len(ac.user) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"user\"")
	}
	if ac.user == nil {
		return nil, errors.New("ent: missing required edge \"user\"")
	}
	return ac.sqlSave(ctx)
}

// SaveX calls Save and panics if Save returns an error.
func (ac *AuthCreate) SaveX(ctx context.Context) *Auth {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ac *AuthCreate) sqlSave(ctx context.Context) (*Auth, error) {
	var (
		res     sql.Result
		builder = sql.Dialect(ac.driver.Dialect())
		a       = &Auth{config: ac.config}
	)
	tx, err := ac.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	insert := builder.Insert(auth.Table).Default()
	if value := ac.created_at; value != nil {
		insert.Set(auth.FieldCreatedAt, *value)
		a.CreatedAt = *value
	}
	if value := ac.updated_at; value != nil {
		insert.Set(auth.FieldUpdatedAt, *value)
		a.UpdatedAt = *value
	}
	if value := ac.access_token; value != nil {
		insert.Set(auth.FieldAccessToken, *value)
		a.AccessToken = *value
	}
	if value := ac.refresh_token; value != nil {
		insert.Set(auth.FieldRefreshToken, *value)
		a.RefreshToken = *value
	}

	id, err := insertLastID(ctx, tx, insert.Returning(auth.FieldID))
	if err != nil {
		return nil, rollback(tx, err)
	}
	a.ID = int(id)
	if len(ac.user) > 0 {
		eid := keys(ac.user)[0]
		query, args := builder.Update(auth.UserTable).
			Set(auth.UserColumn, eid).
			Where(sql.EQ(auth.FieldID, id).And().IsNull(auth.UserColumn)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return nil, rollback(tx, err)
		}
		if int(affected) < len(ac.user) {
			return nil, rollback(tx, &ErrConstraintFailed{msg: fmt.Sprintf("one of \"user\" %v already connected to a different \"Auth\"", keys(ac.user))})
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return a, nil
}