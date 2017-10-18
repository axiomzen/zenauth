package integration

import (
	"context"
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

	ginkgo.Context("GetUsersByIDs", func() {
		var (
			user2, user3 models.User
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

			gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
			statusCode, err = TestRequestV1().
				Post(routes.ResourceUsers + routes.ResourceSignup).
				RequestBody(&signup).
				ResponseBody(&user3).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
		})
		ginkgo.It("Returns the users", func() {
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUsers, err := grpcAuthClient.GetUsersByIDs(ctx, &protobuf.UserIDs{
				Ids: []string{user2.ID, user3.ID},
			})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(grpcUsers.Users)).To(gomega.Equal(2))
			gomega.Expect(grpcUsers.Users[0].Id).To(gomega.Equal(user2.ID))
			gomega.Expect(grpcUsers.Users[0].Email).To(gomega.Equal(user2.Email))
			gomega.Expect(grpcUsers.Users[1].Id).To(gomega.Equal(user3.ID))
			gomega.Expect(grpcUsers.Users[1].Email).To(gomega.Equal(user3.Email))
		})
		ginkgo.It("Returns error if the user is not logged in", func() {
			ctx := getGRPCAuthenticatedContext("INVALID_TOKEN")
			grpcUsers, err := grpcAuthClient.GetUsersByIDs(ctx, &protobuf.UserIDs{
				Ids: []string{user2.ID, user3.ID},
			})
			gomega.Expect(grpcUsers).To(gomega.BeNil())
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("GetUsersByFacebookIDs", func() {
		var (
			user2 models.User
		)

		ginkgo.BeforeEach(func() {
			var fbsignup models.FacebookSignup
			gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
			fbsignup.FacebookID = FacebookTestId
			fbsignup.FacebookToken = FacebookTestToken

			statusCode, err := TestRequestV1().
				Post(routes.ResourceUsers + routes.ResourceFacebookSignup).
				RequestBody(&fbsignup).
				ResponseBody(&user2).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
		})
		ginkgo.AfterEach(func() {
			deleteUser(user2.ID)
		})

		ginkgo.It("Returns the users", func() {
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUsers, err := grpcAuthClient.GetUsersByFacebookIDs(ctx, &protobuf.UserIDs{
				Ids: []string{user2.FacebookID},
			})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(len(grpcUsers.Users)).To(gomega.Equal(1))
			gomega.Expect(grpcUsers.Users[0].Id).To(gomega.Equal(user2.ID))
			gomega.Expect(grpcUsers.Users[0].Email).To(gomega.Equal(user2.Email))
			gomega.Expect(grpcUsers.Users[0].FacebookID).To(gomega.Equal(user2.FacebookID))

		})
		ginkgo.It("Returns error if the user is not logged in", func() {
			ctx := getGRPCAuthenticatedContext("INVALID_TOKEN")
			grpcUsers, err := grpcAuthClient.GetUsersByFacebookIDs(ctx, &protobuf.UserIDs{
				Ids: []string{user2.ID},
			})
			gomega.Expect(grpcUsers).To(gomega.BeNil())
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})

	ginkgo.Context("AuthUserByEmail", func() {
		var (
			signup    protobuf.UserEmailAuth
			protoUser *protobuf.User
			authErr   error
		)
		ginkgo.BeforeEach(func() {
			signup.Email = lorem.Email()
			signup.UserName = lorem.Word(8, 16)
			signup.Password = lorem.Word(8, 16)
		})
		ginkgo.AfterEach(func() {
			deleteUser(protoUser.Id)
		})
		ginkgo.It("Allows user to signup", func() {
			ctx := context.Background()
			protoUser, authErr = grpcAuthClient.AuthUserByEmail(ctx, &signup)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(protoUser.Email).To(gomega.Equal(signup.Email))
			gomega.Expect(protoUser.UserName).To(gomega.Equal(signup.UserName))
			gomega.Expect(protoUser.AuthToken).ToNot(gomega.BeEmpty())

		})
		ginkgo.It("Allows user to login", func() {
			ctx := context.Background()
			// Signup
			_, authErr = grpcAuthClient.AuthUserByEmail(ctx, &signup)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())

			// Login
			protoUser, authErr = grpcAuthClient.AuthUserByEmail(ctx, &signup)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(protoUser.Email).To(gomega.Equal(signup.Email))
			gomega.Expect(protoUser.UserName).To(gomega.Equal(signup.UserName))
			gomega.Expect(protoUser.AuthToken).ToNot(gomega.BeEmpty())
		})
	})

	ginkgo.Context("AuthUserByFacebook", func() {
		var (
			facebook  protobuf.UserFacebookAuth
			protoUser *protobuf.User
			authErr   error
		)
		ginkgo.AfterEach(func() {
			deleteUser(protoUser.Id)
		})
		ginkgo.BeforeEach(func() {
			facebook.FacebookID = FacebookTestId
			facebook.FacebookEmail = lorem.Email()
			facebook.FacebookUsername = lorem.Word(8, 16)
			facebook.FacebookToken = FacebookTestToken
		})
		ginkgo.It("Allows user to signup", func() {
			ctx := context.Background()
			protoUser, authErr = grpcAuthClient.AuthUserByFacebook(ctx, &facebook)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(protoUser.FacebookID).To(gomega.Equal(facebook.FacebookID))
			gomega.Expect(protoUser.AuthToken).ToNot(gomega.BeEmpty())

		})
		ginkgo.It("Allows user to login", func() {
			ctx := context.Background()
			// Signup
			_, authErr := grpcAuthClient.AuthUserByFacebook(ctx, &facebook)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())

			// Login
			protoUser, authErr = grpcAuthClient.AuthUserByFacebook(ctx, &facebook)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(protoUser.FacebookID).To(gomega.Equal(facebook.FacebookID))
			gomega.Expect(protoUser.AuthToken).ToNot(gomega.BeEmpty())
		})

		ginkgo.PIt("Allows repeated usernames by adding number", func() {
			ctx := context.Background()
			signup := protobuf.UserEmailAuth{
				Email:    lorem.Email(),
				UserName: lorem.Word(8, 16),
				Password: lorem.Word(8, 16),
			}

			protoUser, authErr = grpcAuthClient.AuthUserByEmail(ctx, &signup)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(protoUser.UserName).To(gomega.Equal(signup.UserName))

			// Signup
			facebook.FacebookUsername = signup.UserName
			protoUser, authErr = grpcAuthClient.AuthUserByFacebook(ctx, &facebook)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(protoUser.FacebookID).To(gomega.Equal(facebook.FacebookID))
			gomega.Expect(protoUser.UserName).To(gomega.Equal(facebook.FacebookUsername + "-1"))
			gomega.Expect(protoUser.AuthToken).ToNot(gomega.BeEmpty())
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
		ginkgo.It("Allows signin of original account with new linked signin method", func() {
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

			userResponse := res.Users[0]
			// Make the GRPC call
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			_, err = grpcAuthClient.LinkUser(ctx, &protobuf.InvitationCode{
				Type:       constants.InvitationTypeFacebook,
				InviteCode: userResponse.FacebookID,
			})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())

			facebookLogin := models.FacebookUser{
				FacebookID:    FacebookTestId,
				FacebookToken: FacebookTestToken,
			}
			var mergedUser models.User
			statusCode, err = TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebook).
				RequestBody(&facebookLogin).
				ResponseBody(&mergedUser).Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

			// deep compare user to newuser (auth token will be different)
			gomega.Expect(user.ID).To(gomega.Equal(mergedUser.ID))
			gomega.Expect(facebookLogin.FacebookID).To(gomega.Equal(mergedUser.FacebookID))
			gomega.Expect(facebookLogin.FacebookToken).To(gomega.Equal(mergedUser.FacebookToken))

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

	ginkgo.Context("UpdateUserEmail and UpdateUserName", func() {
		var (
			signup    protobuf.UserEmailAuth
			protoUser *protobuf.User
			authErr   error
		)
		ginkgo.BeforeEach(func() {
			signup.Email = lorem.Email()
			signup.UserName = lorem.Word(8, 16)
			signup.Password = lorem.Word(8, 16)
			protoUser, authErr = grpcAuthClient.AuthUserByEmail(context.Background(), &signup)
			gomega.Expect(authErr).ToNot(gomega.HaveOccurred())
		})
		ginkgo.AfterEach(func() {
			deleteUser(protoUser.Id)
		})
		ginkgo.It("Allows user to update email", func() {
			emailUpdate := protobuf.UserEmailAuth{
				Email: lorem.Email(),
			}
			authCtx := getGRPCAuthenticatedContext(protoUser.GetAuthToken())
			updateUser, updateErr := grpcAuthClient.UpdateUserEmail(authCtx, &emailUpdate)

			gomega.Expect(updateErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(updateUser.Email).To(gomega.Equal(emailUpdate.Email))
			gomega.Expect(updateUser.UserName).To(gomega.Equal(protoUser.UserName))
		})
		ginkgo.It("Allows user to update username", func() {
			userNameUpdate := protobuf.UserEmailAuth{
				UserName: lorem.Word(5, 10),
			}
			authCtx := getGRPCAuthenticatedContext(protoUser.GetAuthToken())
			updateUser, updateErr := grpcAuthClient.UpdateUserName(authCtx, &userNameUpdate)

			gomega.Expect(updateErr).ToNot(gomega.HaveOccurred())
			gomega.Expect(updateUser.Email).To(gomega.Equal(protoUser.Email))
			gomega.Expect(updateUser.UserName).To(gomega.Equal(userNameUpdate.UserName))
		})
	})
})
