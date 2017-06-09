package data

import "github.com/axiomzen/zenauth/models"

// GetOrCreateInvitations creates a list of invitations
func (dp *dataProvider) GetOrCreateInvitations(invitations *[]*models.Invitation) error {
	_, err := dp.db.Model(invitations).Create()
	return wrapError(err)
}

// GetInvitationByID Gets an invitation by ID
func (dp *dataProvider) GetInvitationByID(invitation *models.Invitation) error {
	return wrapError(dp.db.Model(invitation).Where("id = ?id").Select())
}
