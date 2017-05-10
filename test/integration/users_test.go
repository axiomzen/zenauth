package integration

import (
	"github.com/axiomzen/authentication/constants"
	"github.com/axiomzen/authentication/helpers"
	"github.com/axiomzen/authentication/models"
	"github.com/axiomzen/authentication/routes"
	"github.com/axiomzen/compare"
	"github.com/axiomzen/golorem"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"net/http"
	"strings"
	"time"

	"fmt"
)

var _ = ginkgo.Describe("Users", func() {

	ginkgo.BeforeEach(func() {
		// clear cache
		// respClearCache := test.PutGatedPage(test.TEST_CACHE).Header(theConf.AuthTokenHeader, TesterToken).Do()
		// gomega.Expect(respClearCache.StatusCode).To(gomega.Equal(http.StatusOK))
	})

	ginkgo.Describe("Signup", func() {

		ginkgo.Context("New user with API token", func() {
			ginkgo.It("should not be able to access GET /VERSION/users/ without signing up", func() {
				statusCode, err := TestRequestV1().Get(routes.ResourceUsers).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})
		})

		ginkgo.Context("User not signed up yet", func() {

			ginkgo.It("should return email does not exist", func() {
				var exists models.Exists
				statusCode, err := TestRequestV1().
					Get(routes.ResourceUsers+routes.ResourceExists).
					URLParam("email", "testuserlive123@gmail.com").
					ResponseBody(&exists).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
				gomega.Expect(exists.Exists).To(gomega.Equal(false))
			})

			ginkgo.It("should be able to sign up and get an authentication token returned that can access pages", func() {
				var signup models.Signup
				gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
				var user models.User
				statusCode, err := TestRequestV1().
					Post(routes.ResourceUsers + routes.ResourceSignup).
					RequestBody(&signup).
					ResponseBody(&user).
					Do()

				defer func() {
					// remove user
					statusCode, err := TestRequestV1().Delete(routes.ResourceTest+routes.ResourceUsers+"/"+user.ID).Header(theConf.AuthTokenHeader, TesterToken).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
				}()

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(user.AuthToken).ToNot(gomega.BeEmpty(), "user.AuthToken")
				gomega.Expect(user.ID).ToNot(gomega.BeEmpty(), "user.ID")
				gomega.Expect(user.Verified).ToNot(gomega.BeTrue())

				// check that the user returned has all the things we expected
				gomega.Expect(signup.FirstName).To(gomega.Equal(user.FirstName))
				gomega.Expect(signup.LastName).To(gomega.Equal(user.LastName))
				gomega.Expect(compare.New().DeepEquals(signup.Email, *user.Email, "signup.Email")).To(gomega.Succeed(), "signup.Email")
				gomega.Expect(user.CreatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Time.IsZero()).NotTo(gomega.BeTrue())

				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

			})

			ginkgo.It("should be able to sign up without a first name or last name", func() {
				var userAuth models.UserAuth
				gomega.Expect(lorem.Fill(&userAuth)).To(gomega.Succeed())
				var userResponse models.User

				statusCode, err := TestRequestV1().
					Post(routes.ResourceUsers + routes.ResourceSignup).
					RequestBody(&struct {
						Email    string `form:"email" json:"email"`
						Password string `form:"password" json:"password"`
					}{*userAuth.Email, userAuth.Password}).
					ResponseBody(&userResponse).
					Do()

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(userResponse.AuthToken).ToNot(gomega.BeEmpty())
				gomega.Expect(userResponse.ID).ToNot(gomega.BeEmpty())

				defer func() {
					// remove user
					statusCode, err := TestRequestV1().Delete(routes.ResourceTest+routes.ResourceUsers+"/"+userResponse.ID).Header(theConf.AuthTokenHeader, TesterToken).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
				}()

				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, userResponse.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

			})

			ginkgo.It("should not be able to sign up without an email", func() {
				var userAuth models.UserAuth
				gomega.Expect(lorem.Fill(&userAuth)).To(gomega.Succeed())
				var errResp models.ErrorResponse

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&struct {
					FirstName string `form:"firstName" json:"firstName"`
					Password  string `form:"password" json:"password"`
				}{*userAuth.FirstName, userAuth.Password}).ErrorResponseBody(&errResp).Do()
				//}{userAuth.FirstName.String, userAuth.Password}).ErrorResponseBody(&errResp).Do()

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				// TODO: refactor all error messages
				gomega.Expect(errResp.GoError.GoError).To(gomega.Equal("Please enter a valid email address"))
			})

			ginkgo.It("should not be able to sign up without a password", func() {
				var userAuth models.UserAuth
				gomega.Expect(lorem.Fill(&userAuth)).To(gomega.Succeed())
				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&struct {
					LastName string `form:"lastName" json:"lastName"`
					Email    string `form:"email" json:"email"`
				}{*userAuth.LastName, *userAuth.Email}).ErrorResponseBody(&errResp).Do()

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				// TODO: brittle test
				gomega.Expect(errResp.GoError.GoError).To(gomega.Equal("Password needs to be at least 8 characters long!"))
			})

			ginkgo.It("should not be able to sign up with a password less than 8 characters", func() {
				var userAuth models.UserAuth
				gomega.Expect(lorem.Fill(&userAuth)).To(gomega.Succeed())
				short := lorem.Word(0, int(theConf.MinPasswordLength)-1)
				userAuth.Password = short
				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&userAuth).ErrorResponseBody(&errResp).Do()

				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				// TODO: brittle test
				gomega.Expect(errResp.GoError.GoError).To(gomega.Equal("Password needs to be at least 8 characters long!"))
			})

		})

		ginkgo.Context("User has signed up", func() {

			var (
				userAuth models.UserAuth
				user     models.User
			)

			ginkgo.BeforeEach(func() {
				gomega.Expect(lorem.Fill(&userAuth)).To(gomega.Succeed())

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).
					RequestBody(&userAuth).ResponseBody(&user).
					//RequestInterceptor(printRequest).
					//ResponseInterceptor(printResponse).
					Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(len(user.ID) > 0).To(gomega.BeTrue())
			})

			ginkgo.AfterEach(func() {
				// delete user
				statusCode, err := TestRequestV1().Delete(routes.ResourceTest+routes.ResourceUsers+"/"+user.ID).Header(theConf.AuthTokenHeader, TesterToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
			})

			ginkgo.It("should not be able to sign up with the same email as an existing user", func() {
				var newUserAuth models.UserAuth
				gomega.Expect(lorem.Fill(&newUserAuth)).To(gomega.Succeed())
				newUserAuth.Email = userAuth.Email
				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&newUserAuth).ErrorResponseBody(&errResp).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusForbidden))
				// TODO: brittle
				//gomega.Expect(errResp.GoError.GoError).To(gomega.Equal("Email already in use/exists"))
				gomega.Expect(errResp.ErrorMessage).To(gomega.Equal("Email already in use/exists"))
			})

			ginkgo.It("should not be able to sign up with the same email (But different case) as an existing user", func() {
				var newUserAuth models.UserAuth
				gomega.Expect(lorem.Fill(&newUserAuth)).To(gomega.Succeed())
				newemail := strings.ToUpper(*userAuth.Email)
				newUserAuth.Email = &newemail
				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&newUserAuth).ErrorResponseBody(&errResp).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusForbidden)) // or StatusBadRequest?
				// TODO: brittle
				gomega.Expect(errResp.ErrorMessage).To(gomega.Equal("Email already in use/exists"))
			})

			ginkgo.It("should return email does exist", func() {
				var exists models.Exists
				statusCode, err := TestRequestV1().
					Get(routes.ResourceUsers+routes.ResourceExists).
					URLParam("email", *userAuth.Email).
					ResponseBody(&exists).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
				gomega.Expect(exists.Exists).To(gomega.Equal(true))
			})

			ginkgo.It("should return email does exist (even with a different case)", func() {
				upperEmail := strings.ToUpper(*userAuth.Email)
				var exists models.Exists
				statusCode, err := TestRequestV1().
					Get(routes.ResourceUsers+routes.ResourceExists).
					URLParam("email", upperEmail).
					ResponseBody(&exists).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
				gomega.Expect(exists.Exists).To(gomega.Equal(true))
			})

			ginkgo.It("should be able to log in with the same credentials as used on sign up and be able to access gated user page", func() {
				auth := models.UserAuth{UserBase: models.UserBase{Email: userAuth.Email}, Password: userAuth.Password}
				var userResponse models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&auth).ResponseBody(&userResponse).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				gomega.Expect(userResponse.AuthToken).ToNot(gomega.BeEmpty())
				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, userResponse.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			})

			ginkgo.It("should return a 440 error when accessing a gated page with expired credentials", func() {
				auth := models.UserAuth{UserBase: models.UserBase{Email: userAuth.Email}, Password: userAuth.Password}
				var userResponse models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&auth).ResponseBody(&userResponse).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				gomega.Expect(userResponse.AuthToken).ToNot(gomega.BeEmpty())
				// use jwt helper
				//fmt.Println("Here hashsecretbytes: " + string(theConf.HashSecretBytes))
				jwt := helpers.JWTHelper{HashSecretBytes: theConf.HashSecretBytes, Token: userResponse.AuthToken}
				gomega.Expect(jwt.Expire()).ToNot(gomega.HaveOccurred())
				userResponse.AuthToken = jwt.Token
				time.Sleep(1 * time.Second)
				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, userResponse.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(constants.StatusExpiredToken))
			})

			ginkgo.It("should be able to log in with the same credentials (But with a different case and whitespace) as used on sign up and be able to access gated user page", func() {
				auth := models.UserAuth{UserBase: models.UserBase{Email: userAuth.Email}, Password: userAuth.Password}
				email := "  " + strings.ToUpper(*auth.Email) + " "
				auth.Email = &email
				var userResponse models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&auth).ResponseBody(&userResponse).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				gomega.Expect(userResponse.AuthToken).ToNot(gomega.BeEmpty())
				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, userResponse.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			})

			ginkgo.It("should not be able to log in with the wrong password", func() {
				login := models.UserAuth{UserBase: models.UserBase{Email: userAuth.Email}, Password: userAuth.Password}
				login.Password = login.Password + "eref"
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})

			ginkgo.It("should not be able to log in with the wrong email", func() {
				login := models.UserAuth{UserBase: models.UserBase{Email: userAuth.Email}, Password: userAuth.Password}
				email := "dfd" + *login.Email
				login.Email = &email
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})

			ginkgo.It("should not be able to log in with no password", func() {
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&struct {
					Email string `form:"email" json:"email"`
				}{*userAuth.Email}).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})

			ginkgo.It("should not be able to log in with no email", func() {
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&struct {
					Password string `form:"password" json:"password"`
				}{userAuth.Password}).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})

			//this test case deals with the "remembered" login details issue
			ginkgo.It("should not be able to log in without a password after a previous login", func() {
				login := models.UserAuth{UserBase: models.UserBase{Email: userAuth.Email}, Password: userAuth.Password}
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				statusCode, err = TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&struct {
					Email string `form:"email" json:"email"`
				}{*userAuth.Email}).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})

			ginkgo.Context("Updating information", func() {

				ginkgo.Context("User has updated all thier information", func() {

					var (
						userUpdate  models.UserUpdate
						updatedUser models.User
						newPassword string
					)

					ginkgo.BeforeEach(func() {
						// update the info
						gomega.Expect(lorem.Fill(&userUpdate)).To(gomega.Succeed())
						userUpdate.ID = user.ID

						// update basic stuff
						statusCode, err := TestRequestV1().Put(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).RequestBody(&userUpdate).ResponseBody(&updatedUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
						// compare stuff
						// because the two structs are so different, we need to manually add them in
						gomega.Expect(userUpdate.Preferences).To(gomega.Equal(updatedUser.Preferences))
						gomega.Expect(userUpdate.FirstName).To(gomega.Equal(updatedUser.FirstName))
						gomega.Expect(userUpdate.LastName).To(gomega.Equal(updatedUser.LastName))
						t1 := time.Unix(updatedUser.UpdatedAt.Time.Unix(), 0)
						gomega.Expect(user.UpdatedAt.Time.Before(updatedUser.UpdatedAt.Time)).To(gomega.BeTrue())

						//gomega.Ω(compare.New().DeepEquals(userUpdate, updatedUser, "userUpdate")).Should(gomega.Succeed())

						// update email
						var userChangeEmail models.UserChangeEmail
						gomega.Expect(lorem.Fill(&userChangeEmail)).To(gomega.Succeed())
						userChangeEmail.ID = user.ID
						statusCode, err = TestRequestV1().Put(routes.ResourceUsers+routes.ResourceEmail).RequestBody(&userChangeEmail).ResponseBody(&updatedUser).Header(theConf.AuthTokenHeader, user.AuthToken).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
						// compare stuff
						// structs are too different (how would we keep track of which anyonymous fields to match etc)
						//gomega.Expect(userChangeEmail.Email).To(gomega.Equal(updatedUser.Email))

						gomega.Ω(compare.New().DeepEquals(userChangeEmail.Email, *updatedUser.Email, "updatedUser.Email")).Should(gomega.Succeed())
						t2 := time.Unix(updatedUser.UpdatedAt.Time.Unix(), 0)
						gomega.Expect(user.UpdatedAt.Time.Before(updatedUser.UpdatedAt.Time)).To(gomega.BeTrue())
						gomega.Expect(t1.Before(updatedUser.UpdatedAt.Time)).To(gomega.BeTrue())

						// statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).Do()

						// gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						// update password
						var userChangePassword models.UserChangePassword
						userChangePassword.OldPassword = userAuth.Password
						// to guaruntee difference
						userChangePassword.NewPassword = userAuth.Password + "$3*"
						newPassword = userChangePassword.NewPassword
						statusCode, err = TestRequestV1().Put(routes.ResourceUsers+routes.ResourcePassword).RequestBody(&userChangePassword).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
						// we can't really compare the password
						gomega.Expect(user.UpdatedAt.Time.Before(updatedUser.UpdatedAt.Time)).To(gomega.BeTrue())
						gomega.Expect(t1.Before(updatedUser.UpdatedAt.Time)).To(gomega.BeTrue())
						gomega.Expect(t2.Before(updatedUser.UpdatedAt.Time)).To(gomega.BeTrue())

						// try a user get
						var userGet models.User
						statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&userGet).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						// compare stuff
						gomega.Ω(compare.New().DeepEquals(updatedUser, userGet, "userGet")).Should(gomega.Succeed())
					})

					ginkgo.It("should be able to login with new credentials after updating them", func() {
						//
						login := models.UserAuth{UserBase: models.UserBase{Email: updatedUser.Email}, Password: newPassword}
						var newUser models.User
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ResponseBody(&newUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						gomega.Expect(newUser.AuthToken).ToNot(gomega.BeEmpty())

						// see if we can log in
						statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
					})

					ginkgo.It("should not be able to login with old credentials after updating them", func() {
						login := models.UserAuth{UserBase: models.UserBase{Email: user.Email}, Password: userAuth.Password}
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
					})

				}) // updated all info

				ginkgo.Context("User has updated some of their info (not email or password)", func() {
					var (
						userUpdate  models.UserUpdate
						updatedUser models.User
					)

					ginkgo.BeforeEach(func() {
						// update the info
						gomega.Expect(lorem.Fill(&userUpdate)).To(gomega.Succeed())

						statusCode, err := TestRequestV1().Put(routes.ResourceUsers).RequestBody(&userUpdate).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						// compare stuff
						gomega.Expect(userUpdate.Preferences).To(gomega.Equal(updatedUser.Preferences))
						gomega.Expect(userUpdate.FirstName).To(gomega.Equal(updatedUser.FirstName))
						gomega.Expect(userUpdate.LastName).To(gomega.Equal(updatedUser.LastName))
						//gomega.Ω(compare.New().DeepEquals(userUpdate, updatedUser, "userUpdate")).Should(gomega.Succeed())

						// do a get, to see if server took the updates
						var newUser models.User
						statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&newUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						// compare stuff
						gomega.Ω(compare.New().DeepEquals(updatedUser, newUser, "updatedUser")).Should(gomega.Succeed())
					})

					ginkgo.It("should still be able to login with same credential if these haven't been updated", func() {
						login := models.UserAuth{UserBase: models.UserBase{Email: user.Email}, Password: userAuth.Password}
						var newUser models.User
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ResponseBody(&newUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						gomega.Expect(newUser.AuthToken).ToNot(gomega.BeEmpty())

						// compare stuff
						// AuthToken will be updated
						gomega.Ω(compare.New().Ignore(".AuthToken").DeepEquals(updatedUser, newUser, "updatedUser")).Should(gomega.Succeed())

						var getUser models.User
						statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, newUser.AuthToken).ResponseBody(&getUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

						// compare stuff
						gomega.Ω(compare.New().Ignore(".AuthToken").DeepEquals(updatedUser, newUser, "updatedUser")).Should(gomega.Succeed())
						// compare stuff
						gomega.Ω(compare.New().DeepEquals(newUser, getUser, "getUser")).Should(gomega.Succeed())
					})
				}) // updated some

				ginkgo.Context("User has updated their password", func() {

					var (
						userChangePassword models.UserChangePassword
					)

					ginkgo.BeforeEach(func() {
						// change the password
						gomega.Expect(lorem.Fill(&userChangePassword)).To(gomega.Succeed())
						userChangePassword.OldPassword = userAuth.Password
						userChangePassword.NewPassword = userChangePassword.OldPassword + "d*7"
						var updatedUser models.User
						statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourcePassword).RequestBody(&userChangePassword).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
					})

					ginkgo.It("should be able to login with the new password", func() {
						var login models.Login
						login.Email = *userAuth.Email
						login.Password = userChangePassword.NewPassword
						var newUser models.User
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ResponseBody(&newUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
						gomega.Expect(newUser.AuthToken).ToNot(gomega.BeEmpty())
						gomega.Expect(newUser.AuthToken).ToNot(gomega.Equal(user.AuthToken))
					})

					ginkgo.It("should not be able to login with their old password", func() {
						var login models.Login
						login.Email = *userAuth.Email
						login.Password = userAuth.Password
						var errResp models.ErrorResponse
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ResponseBody(&errResp).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
					})
				})

				ginkgo.It("should allow a user to update their firstname to blank", func() {
					var userUpdate models.UserUpdate
					gomega.Expect(lorem.Fill(&userUpdate)).To(gomega.Succeed())
					blank := ""
					userUpdate.FirstName = &blank
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers).RequestBody(&userUpdate).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare, first name should be blank
					gomega.Expect(userUpdate.Preferences).To(gomega.Equal(updatedUser.Preferences))
					gomega.Expect(userUpdate.FirstName).To(gomega.Equal(updatedUser.FirstName))
					gomega.Expect(userUpdate.LastName).To(gomega.Equal(updatedUser.LastName))

					//gomega.Ω(compare.New().DeepEquals(update, updatedUser, "update")).Should(gomega.Succeed())

					// do a get, compare result
					var newUser models.User
					statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, updatedUser.AuthToken).ResponseBody(&newUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare stuff
					gomega.Ω(compare.New().DeepEquals(updatedUser, newUser, "updatedUser")).Should(gomega.Succeed())
				})

				ginkgo.It("should allow a user to update their firstname to null", func() {
					var userUpdate models.UserUpdate
					gomega.Expect(lorem.Fill(&userUpdate)).To(gomega.Succeed())
					userUpdate.FirstName = nil
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers).
						RequestBody(&userUpdate).
						Header(theConf.AuthTokenHeader, user.AuthToken).
						ResponseBody(&updatedUser).
						//RequestInterceptor(printRequest).
						//ResponseInterceptor(printResponse).
						Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare, first name should be blank
					gomega.Expect(userUpdate.Preferences).To(gomega.Equal(updatedUser.Preferences))
					gomega.Expect(userUpdate.FirstName).To(gomega.Equal(updatedUser.FirstName))
					gomega.Expect(userUpdate.LastName).To(gomega.Equal(updatedUser.LastName))
					//gomega.Ω(compare.New().DeepEquals(update, updatedUser, "update")).Should(gomega.Succeed())

					// do a get, compare result
					var newUser models.User
					statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, updatedUser.AuthToken).ResponseBody(&newUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare stuff
					gomega.Ω(compare.New().DeepEquals(updatedUser, newUser, "updatedUser")).Should(gomega.Succeed())
				})

				ginkgo.It("should allow a user to update their lastname to blank", func() {
					var userUpdate models.UserUpdate
					gomega.Expect(lorem.Fill(&userUpdate)).To(gomega.Succeed())
					blank := ""
					userUpdate.LastName = &blank
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers).RequestBody(&userUpdate).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare, first name should be blank
					//gomega.Ω(compare.New().DeepEquals(update, updatedUser, "update")).Should(gomega.Succeed())
					gomega.Expect(userUpdate.Preferences).To(gomega.Equal(updatedUser.Preferences))
					gomega.Expect(userUpdate.FirstName).To(gomega.Equal(updatedUser.FirstName))
					gomega.Expect(userUpdate.LastName).To(gomega.Equal(updatedUser.LastName))

					// do a get, compare result
					var newUser models.User
					statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, updatedUser.AuthToken).ResponseBody(&newUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare stuff
					gomega.Ω(compare.New().DeepEquals(updatedUser, newUser, "updatedUser")).Should(gomega.Succeed())
				})

				ginkgo.It("should allow a user to update their lastname to null", func() {
					var userUpdate models.UserUpdate
					gomega.Expect(lorem.Fill(&userUpdate)).To(gomega.Succeed())
					userUpdate.LastName = nil
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers).RequestBody(&userUpdate).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare, first name should be blank
					gomega.Expect(userUpdate.Preferences).To(gomega.Equal(updatedUser.Preferences))
					gomega.Expect(userUpdate.FirstName).To(gomega.Equal(updatedUser.FirstName))
					gomega.Expect(userUpdate.LastName).To(gomega.Equal(updatedUser.LastName))
					//gomega.Ω(compare.New().DeepEquals(update, updatedUser, "update")).Should(gomega.Succeed())

					// do a get, compare result
					var newUser models.User
					statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, updatedUser.AuthToken).ResponseBody(&newUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare stuff
					gomega.Ω(compare.New().DeepEquals(updatedUser, newUser, "updatedUser")).Should(gomega.Succeed())
				})

				ginkgo.It("should not allow a user to pass in blank for email", func() {

					var userChangeEmail models.UserChangeEmail
					userChangeEmail.ID = user.ID
					userChangeEmail.Email = ""
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourceEmail).RequestBody(&userChangeEmail).ResponseBody(&updatedUser).Header(theConf.AuthTokenHeader, user.AuthToken).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))

					// get
					var newUser models.User
					statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&newUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					gomega.Expect(user.Email).To(gomega.Equal(newUser.Email))

				})

				ginkgo.It("should not allow a user to update their password to blank", func() {

					var userChangePassword models.UserChangePassword
					userChangePassword.OldPassword = userAuth.Password
					userChangePassword.NewPassword = ""
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourcePassword).RequestBody(&userChangePassword).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				})

				ginkgo.It("should allow a user to update their password but not change their other info", func() {
					//sleep := exec.Command("sleep", "2")
					//gomega.Expect(sleep.Run()).ToNot(gomega.HaveOccurred())

					var userChangePassword models.UserChangePassword
					gomega.Expect(lorem.Fill(&userChangePassword)).To(gomega.Succeed())
					userChangePassword.OldPassword = userAuth.Password
					userChangePassword.NewPassword = userChangePassword.OldPassword + "d*7"
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourcePassword).RequestBody(&userChangePassword).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// check that nothing else has changed
					fmt.Println("user updatedat: " + user.UpdatedAt.Time.String())
					fmt.Println("updatedUser updatedat: " + updatedUser.UpdatedAt.Time.String())
					gomega.Expect(user.UpdatedAt.Time.Before(updatedUser.UpdatedAt.Time)).To(gomega.BeTrue())

					gomega.Ω(compare.New().Ignore(".UserBase.UpdatedAt").DeepEquals(user, updatedUser, "updatedUser")).Should(gomega.Succeed())

					// Make sure we can login with new credentials
					var newUser models.User
					statusCode, err = TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).
						RequestBody(models.UserAuth{UserBase: models.UserBase{Email: user.Email}, Password: userChangePassword.NewPassword}).
						ResponseBody(&newUser).Do()
					//statusCode, err := makeLoginRequest(`{"email":"testuserlive123@gmail.com", "password":"bobjones"}`)
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// Expect first/last names to stay the same
					gomega.Ω(compare.New().
						IgnoreFields([]string{".UserBase.UpdatedAt", ".AuthToken"}).
						DeepEquals(newUser, user, "updatedUser")).Should(gomega.Succeed())
				})

				ginkgo.It("should not allow a user to update their password to less than the minimum length", func() {
					var userChangePassword models.UserChangePassword
					gomega.Expect(lorem.Fill(&userChangePassword)).To(gomega.Succeed())
					userChangePassword.OldPassword = userAuth.Password
					userChangePassword.NewPassword = "d*7"
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourcePassword).RequestBody(&userChangePassword).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				})

				ginkgo.It("should allow a user to update other fields without sending password or email", func() {
					var userUpdate models.UserUpdate
					gomega.Expect(lorem.Fill(&userUpdate)).To(gomega.Succeed())
					var updatedUser models.User
					statusCode, err := TestRequestV1().Put(routes.ResourceUsers).RequestBody(&struct {
						Firstname string `form:"firstname" json:"firstname"`
						Lastname  string `form:"lastname" json:"lastname"`
					}{*userUpdate.FirstName, *userUpdate.LastName}).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&updatedUser).Do()

					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					// compare to input
					gomega.Expect(userUpdate.FirstName).To(gomega.Equal(updatedUser.FirstName))
					gomega.Expect(userUpdate.LastName).To(gomega.Equal(updatedUser.LastName))

					// compare to old self
					gomega.Ω(compare.New().Ignore(".UserBase.UpdatedAt").
						Ignore(".UserBase.FirstName").
						Ignore(".UserBase.LastName").
						DeepEquals(updatedUser, user, "updatedUser")).Should(gomega.Succeed())

					// json := getJson(respUserInfoPage)

					var newUser models.User
					statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&newUser).Do()
					gomega.Expect(err).ToNot(gomega.HaveOccurred())
					gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

					gomega.Ω(compare.New().Ignore(".UserBase.UpdatedAt").
						Ignore(".UserBase.FirstName").
						Ignore(".UserBase.LastName").
						DeepEquals(newUser, user, "updatedUser")).Should(gomega.Succeed())
				})

				ginkgo.Context("Another user has signed up", func() {

					var (
						anotherSignup models.Signup
						anotherUser   models.User
					)

					ginkgo.BeforeEach(func() {
						gomega.Expect(lorem.Fill(&anotherSignup)).To(gomega.Succeed())
						statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&anotherSignup).ResponseBody(&anotherUser).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
					})

					ginkgo.AfterEach(func() {
						// delete user
						statusCode, err := TestRequestV1().Delete(routes.ResourceTest+routes.ResourceUsers+"/"+anotherUser.ID).Header(theConf.AuthTokenHeader, TesterToken).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
					})

					ginkgo.It("should not be able for another another user to update their email to another users email", func() {
						var changeEmail models.UserChangeEmail
						//gomega.Expect(lorem.Fill(&changeEmail)).To(gomega.Succeed())
						changeEmail.Email = *user.Email
						statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourceEmail).RequestBody(&changeEmail).Header(theConf.AuthTokenHeader, anotherUser.AuthToken).Do()
						gomega.Expect(err).ToNot(gomega.HaveOccurred())
						gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
					})
				})

			})
		})
	})

	ginkgo.Describe("User Login", func() {

		ginkgo.Context("User hasn't signed up yet", func() {
			ginkgo.It("should not be able to log in with blank email and password", func() {
				login := models.Login{}
				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ErrorResponseBody(&errResp).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
				//GoError.GoError
				gomega.Expect(errResp.ErrorMessage).To(gomega.Equal("Invalid email/password combination"))

			})

			ginkgo.It("should not be able to log in with blank email", func() {
				var login models.Login
				gomega.Expect(lorem.Fill(&login)).To(gomega.Succeed())
				login.Email = ""
				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ErrorResponseBody(&errResp).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
				gomega.Expect(errResp.ErrorMessage).To(gomega.Equal("Invalid email/password combination"))
			})

			ginkgo.It("should not be able to log in with blank password", func() {
				var login models.Login
				gomega.Expect(lorem.Fill(&login)).To(gomega.Succeed())
				login.Password = ""
				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ErrorResponseBody(&errResp).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
				gomega.Expect(errResp.ErrorMessage).To(gomega.Equal("Invalid email/password combination"))
			})

			ginkgo.It("should not be able to log in with no input", func() {
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(struct{}{}).Do()
				//gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})
		})

		ginkgo.Context("User has signed up with altering case and whitespace", func() {

			var (
				user      models.User
				signup    models.Signup
				userEmail = "   TestUserLive123@gmail.cOm   "
			)

			ginkgo.BeforeEach(func() {
				gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
				signup.Email = userEmail
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&signup).ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(*user.Email).To(gomega.Equal(strings.ToLower(strings.Trim(signup.Email, " "))))
				//gomega.Expect(compare.New().DeepEquals(*user.Email, ), "user.Email")).To(gomega.Succeed(), "user.Email")
			})

			ginkgo.AfterEach(func() {
				// delete user
				statusCode, err := TestRequestV1().Delete(routes.ResourceTest+routes.ResourceUsers+"/"+user.ID).Header(theConf.AuthTokenHeader, TesterToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
			})

			ginkgo.It("should return email does exist", func() {
				email := strings.ToLower(strings.Trim(userEmail, " "))
				var exists models.Exists
				statusCode, err := TestRequestV1().
					Get(routes.ResourceUsers+routes.ResourceExists).
					URLParam("email", email).
					ResponseBody(&exists).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
				gomega.Expect(exists.Exists).To(gomega.BeTrue())
			})

			ginkgo.It("should return email does exist (even with a different case)", func() {
				email := strings.ToUpper(strings.Trim(userEmail, " "))
				var exists models.Exists
				statusCode, err := TestRequestV1().
					Get(routes.ResourceUsers+routes.ResourceExists).
					URLParam("email", email).
					ResponseBody(&exists).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
				gomega.Expect(exists.Exists).To(gomega.BeTrue())
			})

			ginkgo.It("should be able to log in with the same credentials as used on sign up and be able to access gated user page", func() {
				var login models.Login
				login.Email = *user.Email
				login.Password = signup.Password
				var newUser models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ResponseBody(&newUser).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				gomega.Expect(newUser.AuthToken).ToNot(gomega.BeEmpty())

				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, newUser.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			})

			ginkgo.It("should be able to log in with the same credentials (With different case) as used on sign up and be able to access gated user page", func() {
				var login models.Login
				login.Email = strings.ToUpper(*user.Email)
				login.Password = signup.Password
				var newUser models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceLogin).RequestBody(&login).ResponseBody(&newUser).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				gomega.Expect(newUser.AuthToken).ToNot(gomega.BeEmpty())

				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, newUser.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			})
		})
	})

	ginkgo.Describe("User info page", func() {

		ginkgo.Context("New user/Default state for a user", func() {
			var (
				user   models.User
				signup models.Signup
			)

			ginkgo.BeforeEach(func() {
				gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceSignup).RequestBody(&signup).ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
			})

			ginkgo.AfterEach(func() {
				// delete user
				statusCode, err := TestRequestV1().Delete(routes.ResourceTest+routes.ResourceUsers+"/"+user.ID).Header(theConf.AuthTokenHeader, TesterToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusNoContent))
			})

			ginkgo.It("should return user's id, firstname, lastname, email and empty preferences", func() {

				statusCode, err := TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				gomega.Expect(user.ID).ToNot(gomega.BeEmpty())
				gomega.Expect(user.FirstName).ToNot(gomega.BeNil())
				gomega.Expect(user.LastName).ToNot(gomega.BeNil())
				// TODO: do we want empty preferences or nil?
				gomega.Expect(user.Preferences).To(gomega.BeNil())
				// TODO: more tests around setting and getting preferences (JSON) object
			})
		})
	})
})
