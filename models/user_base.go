package models

import (
	"github.com/axiomzen/null"
)

//go:generate ffjson $GOFILE

// TODO: bson: http://stackoverflow.com/questions/24216510/empty-or-not-required-struct-fields-in-golang (omitempty)
// TODO: pg: http://marcesher.com/2014/10/13/go-working-effectively-with-database-nulls/ `sql:",null"` for structs or sql.NullString etc

// UserBase is the base user struct
type UserBase struct {
	ID        string    `json:"id" lorem:"-" sql:",pk"`
	TableName TableName `json:"-" sql:"users,alias:user" lorem:"-"`

	CreatedAt null.Time `json:"createdAt,omitempty" sql:",null" lorem:"-"`
	UpdatedAt null.Time `json:"updatedAt,omitempty" sql:",null" lorem:"-"`
	Email     string    `json:"email" sql:",null" lorem:"email"`
	UserName  string    `json:"userName" sql:",null" lorem:"uuid"`
	Verified  bool      `json:"verified"`
}
