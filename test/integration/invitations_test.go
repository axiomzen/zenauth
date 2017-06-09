package integration

import (
	"net/http"

	"github.com/axiomzen/golorem"
	"github.com/axiomzen/zenauth/models"
	"github.com/axiomzen/zenauth/protobuf"
	"github.com/axiomzen/zenauth/routes"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Invitations", func() {

	var (
		user *models.User
	)
	ginkgo.BeforeEach(func() {
		var signup models.Signup
		gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())

		statusCode, err := TestRequestV1().
			Post(routes.ResourceUsers + routes.ResourceSignup).
			RequestBody(&signup).
			ResponseBody(&user).
			Do()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
		gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
	})

	ginkgo.Describe("Invite", func() {
		ginkgo.It("can invite a new user by e-mail", func() {
			var res models.InvitationResponse
			req := models.InvitationRequest{
				Emails: []string{"my-friend@zenfriends.com"},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers + routes.ResourceInvitations).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			gomega.Expect(len(res.Users)).To(gomega.Equal(len(req.Emails)))
			gomega.Expect(res.Users[0].Email).To(gomega.Equal(req.Emails[0]))
		})
		ginkgo.It("can fetch an invited user by ID using the users endpoint", func() {
			var res models.InvitationResponse
			req := models.InvitationRequest{
				Emails: []string{"my-friend@zenfriends.com"},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers + routes.ResourceInvitations).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			gomega.Expect(len(res.Users)).To(gomega.Equal(len(req.Emails)))
			var userPublic protobuf.UserPublic
			statusCode, err = TestRequestV1().
				Get(routes.ResourceUsers+"/"+res.Users[0].Id).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				ResponseBody(&userPublic).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

			gomega.Expect(userPublic.Email).To(gomega.Equal(res.Users[0].Email))
			gomega.Expect(userPublic.Id).To(gomega.Equal(res.Users[0].Id))
			gomega.Expect(userPublic.Status).To(gomega.Equal(protobuf.UserStatus_invited))
		})
		ginkgo.It("keeps the same ID after the invited user signs up", func() {})
		ginkgo.It("doesn't invite the same user twice", func() {})
		ginkgo.It("doesn't invite users that already exist", func() {})
		ginkgo.It("fails if there's no token", func() {})
		ginkgo.It("fails if the token is not valid", func() {})
		ginkgo.It("fails if the email is not valid", func() {})
	})
})
