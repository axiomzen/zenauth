package models

// FacebookUpdate for updating the facebook info of a user
type FacebookUpdate struct {
	ID        string    `json:"id" lorem:"-"`
	TableName TableName `sql:"users"       json:"-" lorem:"-"`

	FacebookUser
}
