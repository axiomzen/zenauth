package data

import "github.com/axiomzen/zenauth/models"

// GetOrCreateInvitations creates a list of invitations
func (dp *dataProvider) CreateInvitations(invitations *[]*models.Invitation) error {
	_, err := dp.db.Model(invitations).Create()
	return wrapError(err)
}

// GetInvitationByID Gets an invitation by ID
func (dp *dataProvider) GetInvitationByID(invitation *models.Invitation) error {
	return wrapError(dp.db.Model(invitation).Where("id = ?id").Select())
}

// GetInvitationByEmail gets an invitation by email
func (dp *dataProvider) GetInvitationByEmail(invite *models.Invitation) error {
	return wrapError(dp.db.Model(invite).Where("email = ?email").Select())
}

// DeleteInvitationByEmail deletes the invitation with the email
func (dp *dataProvider) DeleteInvitationByEmail(invite *models.Invitation) error {
	_, err := dp.db.Model(invite).Where("email = ?email").Delete()
	return wrapError(err)
}
