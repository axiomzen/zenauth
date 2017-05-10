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
	TableName TableName `sql:"users,alias:user"       json:"-" lorem:"-"`

	CreatedAt null.Time `bson:"created_at" json:"createdAt,omitempty" sql:",null" lorem:"-"`
	UpdatedAt null.Time `bson:"updated_at" json:"updatedAt,omitempty" sql:",null" lorem:"-"`
	FirstName *string   `bson:"first_name" json:"firstName"  lorem:"word,2,10"`
	LastName  *string   `bson:"last_name" json:"lastName"  lorem:"word,2,10"`
	Email     *string   `bson:"email,omitempty" json:"email"  lorem:"email"`
	Verified  bool      `json:"verified"`
}
