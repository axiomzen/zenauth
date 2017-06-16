package integration

import (
	"net/http"

	"github.com/axiomzen/compare"
	"github.com/axiomzen/golorem"
	"github.com/axiomzen/zenauth/models"
	"github.com/axiomzen/zenauth/routes"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Social Login/Signup Functionality", func() {

	ginkgo.BeforeEach(func() {
		// clear cache (if needed)
		//respClearCache := test.PutGatedPage(test.TEST_CACHE, TesterToken)
		//Expect(respClearCache.StatusCode).To(gomega.Equal(http.StatusOK))
	})

	ginkgo.Describe("Concerning Facebook Signup", func() {

		ginkgo.Context("If user not signed up yet at all", func() {

			ginkgo.It("Should allow Facebook signup with proper token/id", func() {
				//todo
				var fbsignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				fbsignup.FacebookID = FacebookTestId
				fbsignup.FacebookToken = FacebookTestToken

				var user models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).RequestBody(struct {
					FacebookID    string `form:"facebookId"           json:"facebookId"`
					FacebookToken string `form:"facebookToken"        json:"facebookToken"`
					Email         string `form:"email"        json:"email"`
				}{fbsignup.FacebookID, fbsignup.FacebookToken, fbsignup.Email}).ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				defer deleteUser(user.ID)

				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(compare.New().DeepEquals(fbsignup.Email, user.Email, "fbsignup.Email")).To(gomega.Succeed(), "fbsignup.Email")
				//gomega.Expect(compare.New().DeepEquals(fbsignup.FacebookUser, user.FacebookUser, "FacebookUser")).To(gomega.Succeed(), "FacebookUser")
				gomega.Expect(fbsignup.FacebookID).To(gomega.Equal(user.FacebookID))
				gomega.Expect(fbsignup.FacebookToken).To(gomega.Equal(user.FacebookToken))
				gomega.Expect(user.ID).ToNot(gomega.BeEmpty(), "user.ID")
				gomega.Expect(user.AuthToken).ToNot(gomega.BeEmpty(), "user.AuthToken")
				gomega.Expect(user.Verified).ToNot(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
			})

			ginkgo.It("Should allow Facebook signup with proper token/id with no email", func() {
				fbLogin := models.FacebookUser{FacebookID: FacebookTestId, FacebookToken: FacebookTestToken}
				var user models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).RequestBody(&fbLogin).ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				defer deleteUser(user.ID)

				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(user.ID).ToNot(gomega.BeEmpty(), "user.ID")
				gomega.Expect(user.Email).To(gomega.BeEmpty())
				gomega.Expect(fbLogin.FacebookID).To(gomega.Equal(user.FacebookID))
				gomega.Expect(fbLogin.FacebookToken).To(gomega.Equal(user.FacebookToken))
				gomega.Expect(user.AuthToken).ToNot(gomega.BeEmpty(), "user.AuthToken")
				gomega.Expect(user.Verified).ToNot(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
			})

			ginkgo.It("Should allow Facebook signup with proper token/id and save the facebook username", func() {
				var fbsignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				fbsignup.FacebookID = FacebookTestId
				fbsignup.FacebookToken = FacebookTestToken

				var user models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).RequestBody(&fbsignup).ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				defer deleteUser(user.ID)

				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

				// TODO: perhaps sensible deep compare will attempt to locate items
				// found in first struct in the second struct (only err if it couldn't find a match?)
				// substructs must match type?
				// would allow for order independence
				gomega.Expect(compare.New().DeepEquals(fbsignup.FacebookUser, user.FacebookUser, "fbsignup.FacebookUser")).To(gomega.Succeed(), "FacebookUser")

				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(compare.New().DeepEquals(fbsignup.Email, user.Email, "fbsignup.Email")).To(gomega.Succeed(), "fbsignup.Email")
				gomega.Expect(user.ID).ToNot(gomega.BeEmpty(), "user.ID")
				gomega.Expect(user.AuthToken).ToNot(gomega.BeEmpty(), "user.AuthToken")
				gomega.Expect(user.Verified).ToNot(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Time.IsZero()).NotTo(gomega.BeTrue())

			})

			ginkgo.It("Should allow Facebook signup with proper token/id and empty email", func() {
				var fbsignup models.FacebookSignup
				fbsignup.FacebookID = FacebookTestId
				fbsignup.FacebookToken = FacebookTestToken
				fbsignup.Email = ""

				var user models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).RequestBody(&fbsignup).ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())

				defer deleteUser(user.ID)

				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
				gomega.Expect(user.Email).To(gomega.Equal(""))
				gomega.Expect(fbsignup.FacebookID).To(gomega.Equal(user.FacebookID))
				gomega.Expect(fbsignup.FacebookToken).To(gomega.Equal(user.FacebookToken))
				gomega.Expect(user.ID).ToNot(gomega.BeEmpty(), "user.ID")
				gomega.Expect(user.AuthToken).ToNot(gomega.BeEmpty(), "user.AuthToken")
				gomega.Expect(user.Verified).ToNot(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.CreatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Valid).To(gomega.BeTrue())
				gomega.Expect(user.UpdatedAt.Time.IsZero()).NotTo(gomega.BeTrue())
			})

			ginkgo.It("Should not allow Facebook signup with improper token/id", func() {
				// try invalid id
				var fbsignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				// remember generated fake token
				invalidToken := fbsignup.FacebookToken
				fbsignup.FacebookID = "0"
				fbsignup.FacebookToken = FacebookTestToken

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).RequestBody(&fbsignup).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))

				fbsignup.FacebookToken = invalidToken
				fbsignup.FacebookID = FacebookTestId
				statusCode, err = TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).RequestBody(&fbsignup).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
			})
		})

		ginkgo.Context("If user has signed up with Facebook including an email address", func() {

			var (
				user     models.User
				fbsignup models.FacebookSignup
			)

			ginkgo.BeforeEach(func() {

				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				fbsignup.FacebookID = FacebookTestId
				fbsignup.FacebookToken = FacebookTestToken

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).
					RequestBody(&fbsignup).
					ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
			})

			ginkgo.AfterEach(func() {
				deleteUser(user.ID)
			})

			ginkgo.It("Should allow updating email", func() {

				// now add email
				var changeEmail models.UserChangeEmail
				gomega.Expect(lorem.Fill(&changeEmail)).To(gomega.Succeed())

				var updatedUser models.User
				statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourceEmail).
					RequestBody(&changeEmail).
					ResponseBody(&updatedUser).
					Header(theConf.AuthTokenHeader, user.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				// check result
				gomega.Expect(compare.New().IgnoreFields([]string{".UserBase.Email", ".UserBase.UpdatedAt"}).DeepEquals(user, updatedUser, "updatedUser")).To(gomega.Succeed())

				gomega.Expect(changeEmail.Email).To(gomega.Equal(updatedUser.Email))

				var newUser models.User
				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&newUser).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				// check result
				gomega.Expect(compare.New().DeepEquals(updatedUser, newUser, "newUser")).To(gomega.Succeed())
			})

			ginkgo.It("Should not allow signing up with the same fb id", func() {
				var newFbSignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&newFbSignup)).To(gomega.Succeed())
				newFbSignup.FacebookID = FacebookTestId
				newFbSignup.FacebookToken = FacebookTestToken

				var errResponse models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).
					RequestBody(&newFbSignup).
					ErrorResponseBody(&errResponse).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusForbidden))

				// TODO: brittle
				gomega.Expect(errResponse.ErrorMessage).To(gomega.Equal("Social account already exists"))
			})

		})

		ginkgo.Context("If user has signed up with Facebook not including an email address", func() {

			var (
				user models.User
			)

			ginkgo.BeforeEach(func() {

				fbLogin := models.FacebookUser{FacebookID: FacebookTestId, FacebookToken: FacebookTestToken}
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).
					RequestBody(&fbLogin).
					ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
			})

			ginkgo.AfterEach(func() {
				deleteUser(user.ID)
			})

			ginkgo.It("Should still allow updating email", func() {

				// now add email
				var changeEmail models.UserChangeEmail
				gomega.Expect(lorem.Fill(&changeEmail)).To(gomega.Succeed())

				var updatedUser models.User
				statusCode, err := TestRequestV1().Put(routes.ResourceUsers+routes.ResourceEmail).
					RequestBody(&changeEmail).
					ResponseBody(&updatedUser).
					Header(theConf.AuthTokenHeader, user.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				// check result
				gomega.Ω(compare.New().IgnoreFields([]string{".UserBase.Email", ".UserBase.UpdatedAt"}).DeepEquals(user, updatedUser, "updatedUser")).Should(gomega.Succeed())

				gomega.Expect(changeEmail.Email).To(gomega.Equal(updatedUser.Email))

				var newUser models.User
				statusCode, err = TestRequestV1().Get(routes.ResourceUsers).Header(theConf.AuthTokenHeader, user.AuthToken).ResponseBody(&newUser).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				// check result
				gomega.Ω(compare.New().DeepEquals(updatedUser, newUser, "newUser")).Should(gomega.Succeed())
			})

			ginkgo.It("Should not allow signing up with the same fb id", func() {
				var newFbSignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&newFbSignup)).To(gomega.Succeed())
				newFbSignup.FacebookID = FacebookTestId
				newFbSignup.FacebookToken = FacebookTestToken

				var errResponse models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).
					RequestBody(&newFbSignup).
					ErrorResponseBody(&errResponse).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusForbidden))

				// TODO: brittle
				gomega.Expect(errResponse.ErrorMessage).To(gomega.Equal("Social account already exists"))
			})
		})

		ginkgo.Context("If user is signed up already (through email)", func() {
			var (
				user   models.User
				signup models.Signup
			)

			ginkgo.BeforeEach(func() {
				gomega.Expect(lorem.Fill(&signup)).To(gomega.Succeed())
				statusCode, err := TestRequestV1().
					Post(routes.ResourceUsers + routes.ResourceSignup).
					RequestBody(&signup).
					ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))

				gomega.Expect(user.AuthToken).ToNot(gomega.BeEmpty(), "user.AuthToken")
				gomega.Expect(user.ID).ToNot(gomega.BeEmpty(), "user.ID")
				gomega.Expect(user.Verified).ToNot(gomega.BeTrue())
			})

			ginkgo.AfterEach(func() {
				deleteUser(user.ID)
			})

			ginkgo.It("Should Allow Facebook link with proper token/id", func() {
				var fbUpdate models.FacebookUpdate
				gomega.Expect(lorem.Fill(&fbUpdate)).To(gomega.Succeed())
				fbUpdate.FacebookID = FacebookTestId
				fbUpdate.FacebookToken = FacebookTestToken

				// same email?
				var updatedUser models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers+routes.ResourceFacebookLink).
					RequestBody(&fbUpdate).
					ResponseBody(&updatedUser).
					Header(theConf.AuthTokenHeader, user.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				// compare stuff
				gomega.Ω(
					compare.New().
						DeepEquals(fbUpdate.FacebookUser, updatedUser.FacebookUser, "updatedUser")).
					Should(gomega.Succeed())

				gomega.Ω(
					compare.New().
						IgnoreFields([]string{".UserBase.UpdatedAt", ".AuthToken", ".FacebookUser"}).
						DeepEquals(user, updatedUser, "updatedUser")).
					Should(gomega.Succeed())
			})

			ginkgo.It("Should not Allow Facebook link with improper user authtoken", func() {
				var fbsignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				fbsignup.FacebookID = FacebookTestId
				fbsignup.FacebookToken = FacebookTestToken

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers+routes.ResourceFacebookLink).
					RequestBody(&fbsignup).
					Header(theConf.AuthTokenHeader, "garbageAuthtoken").Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusUnauthorized))
			})

			// because I changed functionality here, we should test that
			// attempting to signup again (with same email) should return forbidden
			ginkgo.It("Should Not allow Facebook signup with proper token/id with same email", func() {
				var fbsignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				fbsignup.FacebookID = FacebookTestId
				fbsignup.FacebookToken = FacebookTestToken
				fbsignup.Email = user.Email

				var errResp models.ErrorResponse
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).
					RequestBody(&fbsignup).
					ErrorResponseBody(&errResp).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusForbidden))
				// TODO: Brittle
				gomega.Expect(errResp.ErrorMessage).To(gomega.Equal("Email already in use/exists"))
			})

			ginkgo.It("Should not allow Facebook signup with improper Facebook token/id", func() {
				// user is attempting to facebook signup (even though already signed up)
				// with their existing auth token with invalid stuff
				var fbsignup models.FacebookSignup
				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				// remember generated fake token
				invalidToken := fbsignup.FacebookToken
				fbsignup.FacebookToken = FacebookTestToken

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers+routes.ResourceFacebookSignup).RequestBody(&fbsignup).Header(theConf.AuthTokenHeader, user.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))

				fbsignup.FacebookToken = invalidToken
				fbsignup.FacebookID = FacebookTestId
				statusCode, err = TestRequestV1().Post(routes.ResourceUsers+routes.ResourceFacebookSignup).RequestBody(&fbsignup).Header(theConf.AuthTokenHeader, user.AuthToken).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
			})
		})
	}) // Facebook signup

	ginkgo.Describe("Concerning Facebook Login", func() {

		ginkgo.Context("If user hasn't signed up yet at all", func() {

			ginkgo.It("Should return request for email Facebook login with proper token/id when not signed up", func() {
				fbLogin := models.FacebookUser{FacebookID: FacebookTestId, FacebookToken: FacebookTestToken}
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookLogin).RequestBody(&fbLogin).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusForbidden))

			})

			ginkgo.It("Should not allow Facebook login with improper token/id", func() {
				fbLogin := models.FacebookUser{FacebookID: "0", FacebookToken: FacebookTestToken}
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookLogin).RequestBody(&fbLogin).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))

				fbLogin = models.FacebookUser{FacebookID: FacebookTestId, FacebookToken: "faketoken"}
				statusCode, err = TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookLogin).RequestBody(&fbLogin).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
			})
		})

		ginkgo.Context("If user has signed up via social login", func() {

			var (
				user     models.User
				fbsignup models.FacebookSignup
			)

			ginkgo.BeforeEach(func() {

				gomega.Expect(lorem.Fill(&fbsignup)).To(gomega.Succeed())
				fbsignup.FacebookID = FacebookTestId
				fbsignup.FacebookToken = FacebookTestToken

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookSignup).
					RequestBody(&fbsignup).
					ResponseBody(&user).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusCreated))
			})

			ginkgo.AfterEach(func() {
				deleteUser(user.ID)
			})

			ginkgo.It("Should allow Facebook login with proper token/id", func() {
				var facebookLogin models.FacebookUser
				facebookLogin.FacebookID = fbsignup.FacebookID
				facebookLogin.FacebookToken = fbsignup.FacebookToken
				var newUser models.User
				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookLogin).
					RequestBody(&facebookLogin).
					ResponseBody(&newUser).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))

				// deep compare user to newuser (auth token will be different)
				gomega.Expect(compare.New().Ignore(".AuthToken").DeepEquals(user, newUser, "user")).To(gomega.Succeed(), "FacebookUser")
				gomega.Expect(facebookLogin.FacebookID).To(gomega.Equal(newUser.FacebookID))
				gomega.Expect(facebookLogin.FacebookToken).To(gomega.Equal(newUser.FacebookToken))
			})

			ginkgo.It("Should not allow Facebook login with improper token/id", func() {
				var facebookLogin models.FacebookUser
				facebookLogin.FacebookID = fbsignup.FacebookID
				facebookLogin.FacebookToken = "notarealtoken"

				statusCode, err := TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookLogin).RequestBody(&facebookLogin).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))

				facebookLogin.FacebookID = "0"
				facebookLogin.FacebookToken = fbsignup.FacebookToken

				statusCode, err = TestRequestV1().Post(routes.ResourceUsers + routes.ResourceFacebookLogin).RequestBody(&facebookLogin).Do()
				gomega.Expect(err).ToNot(gomega.HaveOccurred())
				gomega.Expect(statusCode).To(gomega.Equal(http.StatusBadRequest))
			})

		})
	})
})
