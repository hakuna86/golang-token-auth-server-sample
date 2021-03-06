// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"github.com/facebookincubator/ent/dialect/sql/schema"
	"github.com/facebookincubator/ent/schema/field"
)

var (
	// ArticlesColumns holds the columns for the "articles" table.
	ArticlesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "title", Type: field.TypeString},
		{Name: "body", Type: field.TypeString},
		{Name: "star", Type: field.TypeInt},
		{Name: "user_id", Type: field.TypeInt, Nullable: true},
	}
	// ArticlesTable holds the schema information for the "articles" table.
	ArticlesTable = &schema.Table{
		Name:       "articles",
		Columns:    ArticlesColumns,
		PrimaryKey: []*schema.Column{ArticlesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "articles_users_articles",
				Columns: []*schema.Column{ArticlesColumns[6]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// AuthsColumns holds the columns for the "auths" table.
	AuthsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "access_token", Type: field.TypeString},
		{Name: "refresh_token", Type: field.TypeString},
		{Name: "user_id", Type: field.TypeInt, Unique: true, Nullable: true},
	}
	// AuthsTable holds the schema information for the "auths" table.
	AuthsTable = &schema.Table{
		Name:       "auths",
		Columns:    AuthsColumns,
		PrimaryKey: []*schema.Column{AuthsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:  "auths_users_auth",
				Columns: []*schema.Column{AuthsColumns[5]},

				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "updated_at", Type: field.TypeTime},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "username", Type: field.TypeString},
		{Name: "password", Type: field.TypeString},
		{Name: "role", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:        "users",
		Columns:     UsersColumns,
		PrimaryKey:  []*schema.Column{UsersColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ArticlesTable,
		AuthsTable,
		UsersTable,
	}
)

func init() {
	ArticlesTable.ForeignKeys[0].RefTable = UsersTable
	AuthsTable.ForeignKeys[0].RefTable = UsersTable
}
