package integration

import (
	"net/http"
	"time"

	"github.com/axiomzen/golorem"
	"github.com/axiomzen/zenauth/helpers"
	"github.com/axiomzen/zenauth/models"
	"github.com/axiomzen/zenauth/routes"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Password Reset Functionality", func() {

	ginkgo.BeforeEach(func() {
		// clear cache
		//respClearCache := test.PutGatedPage(test.TEST_CACHE, TesterToken)
		//Expect(respClearCache.StatusCode).To(gomega.Equal(http.StatusOK))
	})

	ginkgo.Context("User", func() {
		ginkgo.Context("User has signed up", func() {

			var (
				userAuth models.UserAuth
				user     models.User
			)

			ginkgo.BeforeEach(func() {
				gomega.Expect(lorem.Fill(&userAuth)).To(gomega.Succeed())
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&userAuth).ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
			})

			ginkgo.AfterEach(func() {
				// delete user
				statusCode, err := TestRequestV1().Delete(routes.ResourceTest+routes.ResourceUsers+"/"+user.ID).Header(theConf.AuthTokenHeader, TesterToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
			})

			ginkgo.It("should be able to reset password using correct email", func() {

				// ask for reset
				statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourceForgotPassword).URLParam("email", user.Email).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
			})

			ginkgo.It("should still get status ok if try to reset using an email that is not a valid user", func() {
				// ask for reset
				statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourceForgotPassword).URLParam("email", lorem.Email()).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
			})

			ginkgo.It("should not be able to actually reset their password without any token", func() {
				var userPasswordReset models.UserPasswordReset
				lorem.Fill(&userPasswordReset)
				userPasswordReset.Email = user.Email

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).RequestBody(&userPasswordReset).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
			})

			ginkgo.It("should not be able to actually reset their password without a valid token", func() {
				var userPasswordReset models.UserPasswordReset
				lorem.Fill(&userPasswordReset)
				userPasswordReset.Email = user.Email
				userPasswordReset.Token = "invalidToken"

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).RequestBody(&userPasswordReset).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
			})

			ginkgo.Context("User has successfully requested a reset password token", func() {

				var (
					trt models.TestResetToken
				)

				ginkgo.BeforeEach(func() {

					// hit forget password
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourceForgotPassword).URLParam("email", user.Email).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))

					// use the test route to get the actual token
					statusCode, err = TestRequestV1().
						Get(routes.ResourceTest+routes.ResourceUsers+routes.ResourcePasswordReset).URLParam("email", user.Email).
						ResponseBody(&trt).
						Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
					gomega.Expect(len(trt.Token) > 0).To(gomega.BeTrue())
				})

				ginkgo.It("should not be able to actually reset their password without any token", func() {
					var userPasswordReset models.UserPasswordReset
					lorem.Fill(&userPasswordReset)
					userPasswordReset.Email = user.Email

					statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).RequestBody(&userPasswordReset).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				})

				ginkgo.It("should not be able to actually reset their password without a valid token", func() {
					var userPasswordReset models.UserPasswordReset
					lorem.Fill(&userPasswordReset)
					userPasswordReset.Email = user.Email
					userPasswordReset.Token = "invalidToken"

					statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).RequestBody(&userPasswordReset).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				})

				ginkgo.It("should not be able to actually reset their password with a token not associated with that email", func() {

					var userPasswordReset models.UserPasswordReset
					lorem.Fill(&userPasswordReset)
					userPasswordReset.Email = lorem.Email()
					userPasswordReset.Token = trt.Token

					statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).RequestBody(&userPasswordReset).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				})

				ginkgo.It("should not be able to reset their password to something too short", func() {

					var userPasswordReset models.UserPasswordReset
					lorem.Fill(&userPasswordReset)
					userPasswordReset.Email = user.Email
					userPasswordReset.Token = trt.Token
					userPasswordReset.NewPassword = lorem.Word(0, int(theConf.MinPasswordLength)-1)

					statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).RequestBody(&userPasswordReset).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				})

				// TODO: reset token has expired
				ginkgo.It("should not be able to reset their password with a token that has expired", func() {
					// expire the token
					jwt := helpers.JWTHelper{HashSecretBytes: theConf.HashSecretBytes, Token: user.AuthToken}
					gomega.Expect(jwt.Expire()).ToNot(gomega.HaveOccurred())
					//userResponse.AuthToken = jwt.Token
					time.Sleep(1 * time.Second)

					var userPasswordReset models.UserPasswordReset
					lorem.Fill(&userPasswordReset)
					userPasswordReset.Email = user.Email
					userPasswordReset.Token = jwt.Token

					statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).RequestBody(&userPasswordReset).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				})

				// because this redirects, we can't test this yet
				ginkgo.PIt("should be able to see the status of their reset token", func() {

				})

				ginkgo.Context("User has reset their password", func() {

					var (
						userPasswordReset models.UserPasswordReset
					)

					ginkgo.BeforeEach(func() {
						lorem.Fill(&userPasswordReset)
						userPasswordReset.Email = user.Email
						userPasswordReset.Token = trt.Token

						var newUser models.User
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceResetPassword).
							RequestBody(&userPasswordReset).
							ResponseBody(&newUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
					})

					ginkgo.It("should be able to login with the new password", func() {
						auth := models.UserAuth{UserBase: models.UserBase{Email: user.Email}, Password: userPasswordReset.NewPassword}
						var userResponse models.User
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&auth).ResponseBody(&userResponse).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						gomega.Expect(userResponse.AuthToken).ToNot(gomega.BeEmpty())
						statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, userResponse.AuthToken).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
					})

					ginkgo.It("should be not able to login with the old password", func() {
						auth := models.UserAuth{UserBase: models.UserBase{Email: user.Email}, Password: userAuth.Password}
						var userResponse models.User
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&auth).ResponseBody(&userResponse).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
					})
				})

				// 	//could make a gmail api request but more trouble than it's worth?
				// 	PIt("should have received reset token in email", func() {

				// 	})

				// 	//can't do any of these tests since unable to get token from db using person's email, nor using api to get emailed link
				// 	Context("within reset password token time limit", func() {

				// 		PIt("should be able to get the password reset page", func() {
				// 			// can't do this as the token is only sent in the email
				// 			//respGetTokenPage := test.GetNonGatedPage(test.RESET_USER + "?token=" + resetToken)

				// 			//Expect(respGetTokenPage.StatusCode).To(gomega.Equal(http.StatusOK))
				// 		})

				// 		//doesn't cover actually receiving the token by email
				// 		PIt("should be able to reset password using their emailed token", func() {
				// 			//resetPasswordResp := requestPage("POST", "/reset_user?token="+resetToken, `{"password": "newpassword"}`)
				// 			//Expect(statusCode).To(gomega.Equal(http.StatusOK))
				// 		})

				// 		PIt("should not be possible to reset password if password not long enough", func() {

				// 		})

				// 		PContext("successfully reset password", func() {
				// 			BeforeEach(func() {
				// 				//resetPasswordResp := requestPage("POST", "/reset_user?token="+resetToken, `{"password": "newpassword"}`)
				// 				//Expect(statusCode).To(gomega.Equal(http.StatusOK))
				// 			})
				// 			It("should be able to log in with their new password", func() {
				// 				//userLoginResp := requestPage("POST", LOGIN, `{"email": "testuserlive123@gmail.com", "password": "newpassword"}`)
				// 				//Expect(userLoginResp.StatusCode).To(gomega.Equal(http.StatusOK))
				// 			})

				// 			It("should no longer be able to log in with their old password", func() {
				// 				//userLoginResp := requestPage("POST", LOGIN, `{"email": "testuserlive123@gmail.com", "password": "bobjones"}`)
				// 				//Expect(userLoginResp.StatusCode).To(gomega.Equal(http.StatusOK))
				// 			})

				// 			It("should no longer be able to use the token reset link", func() {
				// 				//resetPasswordResp := requestPage("POST", "/reset_user?token="+resetToken, `{"password": "newpassword"}`)
				// 				//Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				// 			})
				// 		})
			})
		})

		// 	//how to simulate this?
		// 	PContext("outside reset password token time limit", func() {

		// 		It("should not be able to get the password reset page", func() {

		// 		})

		// 		It("should not be able to reset password using the emailed token", func() {

		// 		})
		// 	})
		// })
	})
})
