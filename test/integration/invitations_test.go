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
			email := lorem.Email()
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{email},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			gomega.Expect(len(res.Users)).To(gomega.Equal(len(req.InviteCodes)))
			gomega.Expect(res.Users[0].Email).To(gomega.Equal(req.InviteCodes[0]))
		})
		ginkgo.It("can fetch an invited user by ID using the users endpoint", func() {
			email := lorem.Email()
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{email},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			gomega.Expect(len(res.Users)).To(gomega.Equal(len(req.InviteCodes)))
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
		ginkgo.It("keeps the same ID after the invited user signs up", func() {
			email := lorem.Email()
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{email},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			var userResponse models.User
			var signup models.Signup
			signup.Email = email
			signup.Password = "asdasdasd"
			statusCode, err = TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceSignup).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&signup).
				ResponseBody(&userResponse).
				Do()

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			gomega.Expect(userResponse.ID).To(gomega.Equal(res.Users[0].Id))
		})
		ginkgo.It("doesn't invite the same user twice", func() {
			var email = lorem.Email()
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{email},
			}

			status, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(status).To(gomega.Equal(http.StatusCreated))

			// Try inviting the same email again, should fail
			status, err = TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(status).To(gomega.Equal(http.StatusBadRequest))
		})
		ginkgo.It("doesn't invite users that already exist", func() {
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{user.Email},
			}

			// Try inviting existing user email, should fail
			status, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(status).To(gomega.Equal(http.StatusBadRequest))
		})
		ginkgo.It("fails if there's no token", func() {
			email := lorem.Email()
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{email},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers + routes.ResourceInvitations + routes.ResourceEmail).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
		})
		ginkgo.It("fails if the token is not valid", func() {
			email := lorem.Email()
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{email},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, "definitely a valid token").
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
		})
		ginkgo.It("fails if the email is not valid", func() {
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{"not a valid email at all!!!"},
			}

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceEmail).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
		})
	})
})
