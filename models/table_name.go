package models

//go:generate ffjson $GOFILE

// TableName is a dummy struct to let sql know if the table name
// normally we would just use struct{} for tablename but ffjson
// doesn't like anonymous struct types (TODO: test this assumption)
type TableName struct{}
