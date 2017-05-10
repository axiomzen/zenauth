// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package data

import (
	"errors"
	"github.com/axiomzen/authentication/models"
	"strings"

	"gopkg.in/pg.v4"
	"io/ioutil"
	"regexp"
)

// errNoneAffected there were no rows affected by the call
var errNoneAffected = errors.New("No rows affected!")

// errUniqueEmail returned when the email already exists
var errUniqueEmail = errors.New("Email must be unique")

// wrapError wraps our outgoing error
func wrapError(err error) error {
	if err != nil {
		str := err.Error()
		// if we accidentally passed a DAL error, just let it through
		if _, ok := err.(DALError); ok {
			return err
		}
		if strings.HasPrefix(str, "ERROR #23505") && strings.Contains(str, "users_email_idx") {
			return DALError{Inner: errUniqueEmail, ErrorCode: DALErrorCodeUniqueEmail}
		}
		if strings.HasPrefix(str, "pg: no rows in result set") {
			return DALError{Inner: errNoneAffected, ErrorCode: DALErrorCodeNoneAffected}
		}

		return DALError{Inner: err, ErrorCode: 0}
	}
	return nil
}

// dataProvider the data provider struct
type dataProvider struct {
	db *pg.DB
}

// Ping pings the database to ensure that we can connect to it
func (dp *dataProvider) Ping() (err error) {
	i := 0
	_, err = dp.db.QueryOne(pg.Scan(&i), "SELECT 1")
	return wrapError(err)
}

// closes the database
func (dp *dataProvider) Close() error {
	return wrapError(dp.db.Close())
}

// Create creates the database
func (dp *dataProvider) Create() (err error) {
	defer func() {
		err = wrapError(err)
	}()

	_, err = dp.db.Exec(`DROP DATABASE IF EXISTS "dulpitr9o7a88d"`)
	if err != nil {
		return err
	}

	_, err = dp.db.Exec(`CREATE DATABASE "dulpitr9o7a88d"`)
	return err
}

// Setup sets up the database (adds tables, etc)
func (dp *dataProvider) Setup() (err error) {
	defer func() {
		err = wrapError(err)
	}()

	re, err := regexp.Compile(`(?m)^\s*--.*$`)
	if err != nil {
		return err
	}

	// back out of tests/integration
	data, err := ioutil.ReadFile("../../data/tables.sql")
	if err != nil {
		return err
	}

	// remove comments
	sqlContents := re.ReplaceAllString(string(data), "")

	// try replacing \n with space
	sqlContents = strings.Replace(sqlContents, "\n", " ", -1)

	// debugging
	//fmt.Printf("SQL STATEMENTS: %s", sqlContents)
	_, err = dp.db.Exec(sqlContents)

	return err
}

// Drop removes the database and all data
func (dp *dataProvider) Drop() (err error) {
	_, err = dp.db.Exec(`DROP DATABASE IF EXISTS "dulpitr9o7a88d"`)
	return wrapError(err)
}

// Tx creates a transaction wrapper
func (dp *dataProvider) Tx(fn func(*pg.Tx) error) error {
	return wrapError(dp.db.RunInTransaction(func(tx *pg.Tx) error {
		defer func(t *pg.Tx) {
			if err := recover(); err != nil {
				t.Rollback()
				// rethrow the panic once the database is safe
				panic(err)
			}
		}(tx)
		return fn(tx)
	}))
}

// GetUserByEmail retrieves a user via email
func (dp *dataProvider) GetUserByEmail(user *models.User) error {
	return wrapError(dp.db.Model(user).Where("email = ?email").Select())
}

// GetUserByID retrieves a user via id
func (dp *dataProvider) GetUserByID(user *models.User) error {
	//Where("id = ?id")
	return wrapError(dp.db.Select(user))
}

// UpdateUser updates a user
func (dp *dataProvider) UpdateUser(model interface{}, user *models.User) error {

	res, err := dp.db.Model(model).Returning("*").Update(user)

	if err == nil {
		if res.Affected() != 1 {
			return DALError{Inner: errNoneAffected, ErrorCode: DALErrorCodeNoneAffected}
		}
	}
	return wrapError(err)
}

// UpdateUserVerified will update a users verified field (looking up user by email)
// TODO: think about generating where clause enums
// and 'model/table' enums  so Update.Model(m, data.T).Where(data.T.Y).Returning(&user).Do()
// or do these functions get generated?
func (dp *dataProvider) UpdateUserVerified(user *models.User) error {
	res, err := dp.db.Model(user).Set("verified = ?verified").Where("email = ?email").Returning("*").Update()
	if err == nil {
		if res.Affected() != 1 {
			return DALError{Inner: errNoneAffected, ErrorCode: DALErrorCodeNoneAffected}
		}
	}
	return wrapError(err)
}

// CreateUserResetToken will update a users password reset token based on email
func (dp *dataProvider) CreateUserResetToken(user *models.User) error {
	res, err := dp.db.Model(user).Set("reset_token = ?reset_token").Where("email = ?email").Returning("*").Update()
	if err == nil {
		if res.Affected() != 1 {
			return DALError{Inner: errNoneAffected, ErrorCode: DALErrorCodeNoneAffected}
		}
	}
	return wrapError(err)
}

// ConsumeUserResetToken will do a bunch of stuff
func (dp *dataProvider) ConsumeUserResetToken(user *models.User) error {
	res, err := dp.db.Model(user).Set("reset_token = NULL, hash = ?hash").Where("email = ?email AND reset_token = ?reset_token").Returning("*").Update()
	if err == nil {
		if res.Affected() != 1 {
			return DALError{Inner: errNoneAffected, ErrorCode: DALErrorCodeNoneAffected}
		}
	}
	return wrapError(err)
}

// CreateUser creates a user
func (dp *dataProvider) CreateUser(user *models.User) error {
	return wrapError(dp.db.Create(user))
}

// DeleteUser deletes a user (by user id)
func (dp *dataProvider) DeleteUser(user *models.User) error {
	return wrapError(dp.db.Delete(user))
}

// ChangeUserPassword allows you to change the password of a user
func (dp *dataProvider) UpdateUserHash(newHash string, user *models.User) error {
	// TODO: we need to err if no rows were updated
	res, err := dp.db.Model(user).Set("hash = ?", newHash).Where("id = ?id AND hash = ?hash").Returning("*").Update()
	if err == nil {
		if res.Affected() != 1 {
			return DALError{Inner: errNoneAffected, ErrorCode: DALErrorCodeNoneAffected}
		}
	}
	return wrapError(err)
}

// ClearUserResetToken sets the reset token to nil (test route)
func (dp *dataProvider) ClearUserResetToken(user *models.User) error {
	res, err := dp.db.Model(user).Set("reset_token = ?reset_token").Where("id = ?id").Returning("*").Update()
	if err == nil {
		if res.Affected() != 1 {
			return DALError{Inner: errNoneAffected, ErrorCode: DALErrorCodeNoneAffected}
		}
	}
	return wrapError(err)
}
