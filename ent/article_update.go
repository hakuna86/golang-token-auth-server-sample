// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/hakuna86/golang-token-auth-server-sample/ent/article"
	"github.com/hakuna86/golang-token-auth-server-sample/ent/predicate"
	"github.com/hakuna86/golang-token-auth-server-sample/ent/user"
)

// ArticleUpdate is the builder for updating Article entities.
type ArticleUpdate struct {
	config

	updated_at  *time.Time
	title       *string
	body        *string
	star        *int
	addstar     *int
	user        map[int]struct{}
	clearedUser bool
	predicates  []predicate.Article
}

// Where adds a new predicate for the builder.
func (au *ArticleUpdate) Where(ps ...predicate.Article) *ArticleUpdate {
	au.predicates = append(au.predicates, ps...)
	return au
}

// SetUpdatedAt sets the updated_at field.
func (au *ArticleUpdate) SetUpdatedAt(t time.Time) *ArticleUpdate {
	au.updated_at = &t
	return au
}

// SetTitle sets the title field.
func (au *ArticleUpdate) SetTitle(s string) *ArticleUpdate {
	au.title = &s
	return au
}

// SetBody sets the body field.
func (au *ArticleUpdate) SetBody(s string) *ArticleUpdate {
	au.body = &s
	return au
}

// SetStar sets the star field.
func (au *ArticleUpdate) SetStar(i int) *ArticleUpdate {
	au.star = &i
	au.addstar = nil
	return au
}

// AddStar adds i to star.
func (au *ArticleUpdate) AddStar(i int) *ArticleUpdate {
	if au.addstar == nil {
		au.addstar = &i
	} else {
		*au.addstar += i
	}
	return au
}

// SetUserID sets the user edge to User by id.
func (au *ArticleUpdate) SetUserID(id int) *ArticleUpdate {
	if au.user == nil {
		au.user = make(map[int]struct{})
	}
	au.user[id] = struct{}{}
	return au
}

// SetNillableUserID sets the user edge to User by id if the given value is not nil.
func (au *ArticleUpdate) SetNillableUserID(id *int) *ArticleUpdate {
	if id != nil {
		au = au.SetUserID(*id)
	}
	return au
}

// SetUser sets the user edge to User.
func (au *ArticleUpdate) SetUser(u *User) *ArticleUpdate {
	return au.SetUserID(u.ID)
}

// ClearUser clears the user edge to User.
func (au *ArticleUpdate) ClearUser() *ArticleUpdate {
	au.clearedUser = true
	return au
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (au *ArticleUpdate) Save(ctx context.Context) (int, error) {
	if au.updated_at == nil {
		v := article.UpdateDefaultUpdatedAt()
		au.updated_at = &v
	}
	if len(au.user) > 1 {
		return 0, errors.New("ent: multiple assignments on a unique edge \"user\"")
	}
	return au.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (au *ArticleUpdate) SaveX(ctx context.Context) int {
	affected, err := au.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (au *ArticleUpdate) Exec(ctx context.Context) error {
	_, err := au.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (au *ArticleUpdate) ExecX(ctx context.Context) {
	if err := au.Exec(ctx); err != nil {
		panic(err)
	}
}

func (au *ArticleUpdate) sqlSave(ctx context.Context) (n int, err error) {
	var (
		builder  = sql.Dialect(au.driver.Dialect())
		selector = builder.Select(article.FieldID).From(builder.Table(article.Table))
	)
	for _, p := range au.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = au.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := au.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(article.Table)
	)
	updater = updater.Where(sql.InInts(article.FieldID, ids...))
	if value := au.updated_at; value != nil {
		updater.Set(article.FieldUpdatedAt, *value)
	}
	if value := au.title; value != nil {
		updater.Set(article.FieldTitle, *value)
	}
	if value := au.body; value != nil {
		updater.Set(article.FieldBody, *value)
	}
	if value := au.star; value != nil {
		updater.Set(article.FieldStar, *value)
	}
	if value := au.addstar; value != nil {
		updater.Add(article.FieldStar, *value)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if au.clearedUser {
		query, args := builder.Update(article.UserTable).
			SetNull(article.UserColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if len(au.user) > 0 {
		for eid := range au.user {
			query, args := builder.Update(article.UserTable).
				Set(article.UserColumn, eid).
				Where(sql.InInts(article.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return 0, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
}

// ArticleUpdateOne is the builder for updating a single Article entity.
type ArticleUpdateOne struct {
	config
	id int

	updated_at  *time.Time
	title       *string
	body        *string
	star        *int
	addstar     *int
	user        map[int]struct{}
	clearedUser bool
}

// SetUpdatedAt sets the updated_at field.
func (auo *ArticleUpdateOne) SetUpdatedAt(t time.Time) *ArticleUpdateOne {
	auo.updated_at = &t
	return auo
}

// SetTitle sets the title field.
func (auo *ArticleUpdateOne) SetTitle(s string) *ArticleUpdateOne {
	auo.title = &s
	return auo
}

// SetBody sets the body field.
func (auo *ArticleUpdateOne) SetBody(s string) *ArticleUpdateOne {
	auo.body = &s
	return auo
}

// SetStar sets the star field.
func (auo *ArticleUpdateOne) SetStar(i int) *ArticleUpdateOne {
	auo.star = &i
	auo.addstar = nil
	return auo
}

// AddStar adds i to star.
func (auo *ArticleUpdateOne) AddStar(i int) *ArticleUpdateOne {
	if auo.addstar == nil {
		auo.addstar = &i
	} else {
		*auo.addstar += i
	}
	return auo
}

// SetUserID sets the user edge to User by id.
func (auo *ArticleUpdateOne) SetUserID(id int) *ArticleUpdateOne {
	if auo.user == nil {
		auo.user = make(map[int]struct{})
	}
	auo.user[id] = struct{}{}
	return auo
}

// SetNillableUserID sets the user edge to User by id if the given value is not nil.
func (auo *ArticleUpdateOne) SetNillableUserID(id *int) *ArticleUpdateOne {
	if id != nil {
		auo = auo.SetUserID(*id)
	}
	return auo
}

// SetUser sets the user edge to User.
func (auo *ArticleUpdateOne) SetUser(u *User) *ArticleUpdateOne {
	return auo.SetUserID(u.ID)
}

// ClearUser clears the user edge to User.
func (auo *ArticleUpdateOne) ClearUser() *ArticleUpdateOne {
	auo.clearedUser = true
	return auo
}

// Save executes the query and returns the updated entity.
func (auo *ArticleUpdateOne) Save(ctx context.Context) (*Article, error) {
	if auo.updated_at == nil {
		v := article.UpdateDefaultUpdatedAt()
		auo.updated_at = &v
	}
	if len(auo.user) > 1 {
		return nil, errors.New("ent: multiple assignments on a unique edge \"user\"")
	}
	return auo.sqlSave(ctx)
}

// SaveX is like Save, but panics if an error occurs.
func (auo *ArticleUpdateOne) SaveX(ctx context.Context) *Article {
	a, err := auo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return a
}

// Exec executes the query on the entity.
func (auo *ArticleUpdateOne) Exec(ctx context.Context) error {
	_, err := auo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (auo *ArticleUpdateOne) ExecX(ctx context.Context) {
	if err := auo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (auo *ArticleUpdateOne) sqlSave(ctx context.Context) (a *Article, err error) {
	var (
		builder  = sql.Dialect(auo.driver.Dialect())
		selector = builder.Select(article.Columns...).From(builder.Table(article.Table))
	)
	article.ID(auo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = auo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		a = &Article{config: auo.config}
		if err := a.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into Article: %v", err)
		}
		id = a.ID
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("Article with id: %v", auo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one Article with the same id: %v", auo.id)
	}

	tx, err := auo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(article.Table)
	)
	updater = updater.Where(sql.InInts(article.FieldID, ids...))
	if value := auo.updated_at; value != nil {
		updater.Set(article.FieldUpdatedAt, *value)
		a.UpdatedAt = *value
	}
	if value := auo.title; value != nil {
		updater.Set(article.FieldTitle, *value)
		a.Title = *value
	}
	if value := auo.body; value != nil {
		updater.Set(article.FieldBody, *value)
		a.Body = *value
	}
	if value := auo.star; value != nil {
		updater.Set(article.FieldStar, *value)
		a.Star = *value
	}
	if value := auo.addstar; value != nil {
		updater.Add(article.FieldStar, *value)
		a.Star += *value
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if auo.clearedUser {
		query, args := builder.Update(article.UserTable).
			SetNull(article.UserColumn).
			Where(sql.InInts(user.FieldID, ids...)).
			Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if len(auo.user) > 0 {
		for eid := range auo.user {
			query, args := builder.Update(article.UserTable).
				Set(article.UserColumn, eid).
				Where(sql.InInts(article.FieldID, ids...)).
				Query()
			if err := tx.Exec(ctx, query, args, &res); err != nil {
				return nil, rollback(tx, err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return a, nil
}
