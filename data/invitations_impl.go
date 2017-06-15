package data

import (
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/models"
)

// CreateInvitations creates a list of invitations
func (dp *dataProvider) CreateInvitations(invitations *models.Invitations) error {
	_, err := dp.db.Model(invitations).Create()
	return wrapError(err)
}

// GetInvitationByID Gets an invitation by ID
func (dp *dataProvider) GetInvitationByID(invitation *models.Invitation) error {
	return wrapError(dp.db.Model(invitation).Where("id = ?id").Select())
}

// GetAllInvitations Gets all invitations
func (dp *dataProvider) GetAllInvitations(invitations *models.Invitations) error {
	return wrapError(dp.db.Model(invitations).Select())
}

// GetInvitationByEmail gets an invitation by email
func (dp *dataProvider) GetInvitationByEmail(invite *models.Invitation) error {
	return wrapError(dp.db.Model(invite).Where("type = ?", constants.InvitationTypeEmail).Where("code = ?code").Select())
}

// DeleteInvitationByEmail deletes the invitation with the email
func (dp *dataProvider) DeleteInvitationByEmail(invite *models.Invitation) error {
	_, err := dp.db.Model(invite).Where("type = ?", constants.InvitationTypeEmail).Where("code = ?code").Delete()
	return wrapError(err)
}

// GetInvitation gets an invitation based on Type field
func (dp *dataProvider) GetInvitation(invite *models.Invitation) error {
	return wrapError(dp.db.Model(invite).Where("type = ?type").Where("code = ?code").Select())
}

// DeleteInvitation deletes the invitation based on Type field
func (dp *dataProvider) DeleteInvitation(invite *models.Invitation) error {
	_, err := dp.db.Model(invite).Where("type = ?type").Where("code = ?code").Delete()
	return wrapError(err)
}
