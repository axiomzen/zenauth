package integration

import (
	"net/http"

	lorem "github.com/axiomzen/golorem"
	"github.com/axiomzen/zenauth/constants"
	"github.com/axiomzen/zenauth/models"
	"github.com/axiomzen/zenauth/protobuf"
	"github.com/axiomzen/zenauth/routes"
	pEmpty "github.com/golang/protobuf/ptypes/empty"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Auth GRPC", func() {
	var (
		user models.User
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

	ginkgo.AfterEach(func() {
		deleteUser(user.ID)
	})

	ginkgo.Context("GetCurrentUser", func() {
		ginkgo.It("Returns the user", func() {
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUser, err := grpcAuthClient.GetCurrentUser(ctx, &pEmpty.Empty{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(grpcUser.Id).To(gomega.Equal(user.ID))
			gomega.Expect(grpcUser.Email).To(gomega.Equal(user.Email))
		})
		ginkgo.It("Returns error if the user is not logged in", func() {
			ctx := getGRPCAuthenticatedContext("INVALID_TOKEN")
			grpcUser, err := grpcAuthClient.GetCurrentUser(ctx, &pEmpty.Empty{})
			gomega.Expect(grpcUser).To(gomega.BeNil())
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("GetUserByID", func() {
		var (
			user2 models.User
		)

		ginkgo.BeforeEach(func() {
			var signup models.Signup
			gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers + routes.ResourceSignup).
				RequestBody(&signup).
				ResponseBody(&user2).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
		})
		ginkgo.It("Returns the user", func() {
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUser, err := grpcAuthClient.GetUserByID(ctx, &protobuf.UserID{
				Id: user2.ID,
			})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(grpcUser.Id).To(gomega.Equal(user2.ID))
			gomega.Expect(grpcUser.Email).To(gomega.Equal(user2.Email))
		})
		ginkgo.It("Returns error if the user is not logged in", func() {
			ctx := getGRPCAuthenticatedContext("INVALID_TOKEN")
			grpcUser, err := grpcAuthClient.GetUserByID(ctx, &protobuf.UserID{
				Id: user2.ID,
			})
			gomega.Expect(grpcUser).To(gomega.BeNil())
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("LinkUser", func() {

		ginkgo.It("Returns the user to be merged", func() {
			// Create the facebook user
			var userResponse models.User
			signup := models.FacebookSignup{
				FacebookUser: models.FacebookUser{
					FacebookID:    FacebookTestId,
					FacebookToken: FacebookTestToken,
				},
			}
			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers + routes.ResourceFacebook).
				RequestBody(&signup).
				ResponseBody(&userResponse).
				Do()

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			// Make the GRPC call
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUser, err := grpcAuthClient.LinkUser(ctx, &protobuf.InvitationCode{
				Type:       constants.InvitationTypeFacebook,
				InviteCode: userResponse.FacebookID,
			})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(grpcUser.Id).To(gomega.Equal(userResponse.ID))
			gomega.Expect(grpcUser.FacebookID).To(gomega.Equal(userResponse.FacebookID))
			gomega.Expect(grpcUser.Status).To(gomega.Equal(protobuf.UserStatus_merged))
		})
		ginkgo.It("Returns the invite to be merged", func() {
			// Invite Test FB user
			var res models.InvitationResponse
			req := models.InvitationRequest{
				InviteCodes: []string{FacebookTestId},
			}
			defer clearInvitations()

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers+routes.ResourceInvitations+routes.ResourceFacebook).
				Header(theConf.AuthTokenHeader, user.AuthToken).
				RequestBody(&req).
				ResponseBody(&res).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

			gomega.Expect(len(res.Users)).To(gomega.Equal(len(req.InviteCodes)))
			userResponse := res.Users[0]
			gomega.Expect(res.Users[0].FacebookID).To(gomega.Equal(req.InviteCodes[0]))

			// Make the GRPC call
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUser, err := grpcAuthClient.LinkUser(ctx, &protobuf.InvitationCode{
				Type:       constants.InvitationTypeFacebook,
				InviteCode: userResponse.FacebookID,
			})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(grpcUser.Id).To(gomega.Equal(userResponse.Id))
			gomega.Expect(grpcUser.FacebookID).To(gomega.Equal(userResponse.FacebookID))
			gomega.Expect(grpcUser.Status).To(gomega.Equal(protobuf.UserStatus_merged))
		})
		ginkgo.It("Returns original user if no invite/user associated with request", func() {
			fbid := "new_facebook_id"
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUser, err := grpcAuthClient.LinkUser(ctx, &protobuf.InvitationCode{
				Type:       constants.InvitationTypeFacebook,
				InviteCode: fbid,
			})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(grpcUser.Id).To(gomega.Equal(user.ID))
			gomega.Expect(grpcUser.FacebookID).To(gomega.Equal(fbid))
			gomega.Expect(grpcUser.Status).To(gomega.Equal(protobuf.UserStatus_created))
		})
		ginkgo.It("Returns error if the invite type is not accepted", func() {
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUser, err := grpcAuthClient.LinkUser(ctx, &protobuf.InvitationCode{
				Type:       "INVALID_TYPE",
				InviteCode: lorem.Word(10, 20),
			})
			gomega.Expect(grpcUser).To(gomega.BeNil())
			gomega.Expect(err).To(gomega.HaveOccurred())
		})

	})
})
