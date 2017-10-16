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
	// GetUserByUserName retrieves a user via username
	GetUserByUserName(user *models.User) error
	// GetUserByEmailOrUserName retrieves a user via email or username
	GetUserByEmailOrUserName(user *models.User) error
	// GetUser retrieves a user via id
	GetUserByID(user *models.User) error
	// GetUsersByIDs retrieves users by their ids
	GetUsersByIDs(users *models.Users) error
	// GetUserByFacebookID retrieves a user from the facebook id
	GetUserByFacebookID(user *models.User) error
	// GetUsersByFacebookIDs gets a list of users by facebook ids
	GetUsersByFacebookIDs(ids []string, users *models.Users) error
	// UpdateUserFacebookInfo updates the user's facebook token
	UpdateUserFacebookInfo(user *models.User) error
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
	// MergeUsers merges the users, with the first user taking precedence.
	MergeUsers(firstUser, secondUser *models.User) error
	// GetUsernameCount counts this username
	GetUsernameCount(username string) (int, error)

	// CreateInvitations creates a list of invitations
	CreateInvitations(invitations *models.Invitations) error
	// GetInvitation gets an invite by type and invite code
	GetInvitation(invite *models.Invitation) error
	// GetAllInvitations gets all invitations
	GetAllInvitations(invitations *models.Invitations) error
	// DeleteInvitation deletes an invite by type and invite code
	DeleteInvitation(invite *models.Invitation) error
	// GetInvitationByID Gets an invitation by ID
	GetInvitationByID(invitation *models.Invitation) error
	// GetInvitationByEmail gets an invite by email
	GetInvitationByEmail(invite *models.Invitation) error
	// DeleteInvitationByEmail deletes an invite by email
	DeleteInvitationByEmail(invite *models.Invitation) error
}
