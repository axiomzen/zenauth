// THIS FILE WAS HATCHED WITH github.com/axiomzen/hatch
// THIS FILE IS SAFE TO EDIT BY HAND

package data

import (
	"github.com/axiomzen/zenauth/models"

	pg "gopkg.in/pg.v4"
)

// Provider the interface all data providers must implement at minimum
type Provider interface {
	// Ping allows one to test the connectivity of the DB
	Ping() error
	// Close closes all connections to the database
	Close() error
	// Create creates the database
	Create() error
	// Setup sets up the database (adds tables, etc)
	Setup() error
	// Drop removes the database and all data
	Drop() error

	// Tx opens a transaction wrapper
	Tx(func(*pg.Tx) error) error
}

// ZENAUTHProvider is the data provider for this app
// TODO: consistency in API usage (pass in struct or return it?)
type ZENAUTHProvider interface {
	Provider

	// GetUserByEmail retrieves a user via email
	GetUserByEmail(user *models.User) error
	// GetUser retrieves a user via id
	GetUserByID(user *models.User) error
	// GetUserByResetToken returns the user via reset token
	//GetUserByResetToken(resetToken string, user *models.User) error

	// UpdateUser updates a user (takes in interface because we want to accept all updates eventually)
	UpdateUser(update interface{}, user *models.User) error
	// UpdateUserVerified will update a users verified field (looking up user by email)
	UpdateUserVerified(user *models.User) error
	// UpdateUserHash allows you to change the password of a user
	UpdateUserHash(newHash string, user *models.User) error
	// CreateUserResetToken will update a users reset token based on email
	CreateUserResetToken(user *models.User) error
	// ConsumeUserResetToken will do a bunch of stuff
	ConsumeUserResetToken(user *models.User) error
	// ClearUserResetToken
	ClearUserResetToken(user *models.User) error
	// CreateUser creates a user
	CreateUser(user *models.User) error
	// DeleteUser deletes a user (by user id)
	DeleteUser(user *models.User) error

	// GetOrCreateInvitations creates a list of invitations
	GetOrCreateInvitations(invitations *[]*models.Invitation) error
	// GetInvitationByID Gets an invitation by ID
	GetInvitationByID(invitation *models.Invitation) error
}
