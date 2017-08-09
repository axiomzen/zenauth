package models

//go:generate ffjson $GOFILE

// UserChangeEmail everything you need for changing email, nothing you don't
type UserChangeEmail struct {
	TableName TableName `sql:"users"       json:"-" lorem:"-"`
	ID        string    `bson:"id" json:"id" lorem:"-"`

	Email string `bson:"email"               json:"email"  lorem:"email"`
}

// UserChangeUserName everything you need for changing email, nothing you don't
type UserChangeUserName struct {
	TableName TableName `sql:"users"       json:"-" lorem:"-"`
	ID        string    `json:"id" lorem:"-"`

	UserName string `json:"userName"  lorem:"word,5,10"`
}
