// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// Auth is the model entity for the Auth schema.
type Auth struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// AccessToken holds the value of the "access_token" field.
	AccessToken string `json:"access_token,omitempty"`
	// RefreshToken holds the value of the "refresh_token" field.
	RefreshToken string `json:"refresh_token,omitempty"`
}

// FromRows scans the sql response data into Auth.
func (a *Auth) FromRows(rows *sql.Rows) error {
	var scana struct {
		ID           int
		CreatedAt    sql.NullTime
		UpdatedAt    sql.NullTime
		AccessToken  sql.NullString
		RefreshToken sql.NullString
	}
	// the order here should be the same as in the `auth.Columns`.
	if err := rows.Scan(
		&scana.ID,
		&scana.CreatedAt,
		&scana.UpdatedAt,
		&scana.AccessToken,
		&scana.RefreshToken,
	); err != nil {
		return err
	}
	a.ID = scana.ID
	a.CreatedAt = scana.CreatedAt.Time
	a.UpdatedAt = scana.UpdatedAt.Time
	a.AccessToken = scana.AccessToken.String
	a.RefreshToken = scana.RefreshToken.String
	return nil
}

// QueryUser queries the user edge of the Auth.
func (a *Auth) QueryUser() *UserQuery {
	return (&AuthClient{a.config}).QueryUser(a)
}

// Update returns a builder for updating this Auth.
// Note that, you need to call Auth.Unwrap() before calling this method, if this Auth
// was returned from a transaction, and the transaction was committed or rolled back.
func (a *Auth) Update() *AuthUpdateOne {
	return (&AuthClient{a.config}).UpdateOne(a)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (a *Auth) Unwrap() *Auth {
	tx, ok := a.config.driver.(*txDriver)
	if !ok {
		panic("ent: Auth is not a transactional entity")
	}
	a.config.driver = tx.drv
	return a
}

// String implements the fmt.Stringer.
func (a *Auth) String() string {
	var builder strings.Builder
	builder.WriteString("Auth(")
	builder.WriteString(fmt.Sprintf("id=%v", a.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(a.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(a.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", access_token=")
	builder.WriteString(a.AccessToken)
	builder.WriteString(", refresh_token=")
	builder.WriteString(a.RefreshToken)
	builder.WriteByte(')')
	return builder.String()
}

// Auths is a parsable slice of Auth.
type Auths []*Auth

// FromRows scans the sql response data into Auths.
func (a *Auths) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scana := &Auth{}
		if err := scana.FromRows(rows); err != nil {
			return err
		}
		*a = append(*a, scana)
	}
	return nil
}

func (a Auths) config(cfg config) {
	for _i := range a {
		a[_i].config = cfg
	}
}
