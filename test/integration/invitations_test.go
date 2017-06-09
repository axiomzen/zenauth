package integration

import (
	"net/http"

	"github.com/axiomzen/golorem"
	"github.com/axiomzen/zenauth/models"
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
		ginkgo.FIt("can invite a new user by e-mail", func() {
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
	})
})
