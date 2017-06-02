package integration

import (
	"net/http"

	lorem "github.com/axiomzen/golorem"
	"github.com/axiomzen/zenauth/models"
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

	ginkgo.Context("GetUser", func() {
		ginkgo.It("Returns the user", func() {
			ctx := getGRPCAuthenticatedContext(user.AuthToken)
			grpcUser, err := grpcAuthClient.GetUser(ctx, &pEmpty.Empty{})
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(grpcUser.Id).To(gomega.Equal(user.ID))
			gomega.Expect(grpcUser.Email).To(gomega.Equal(user.Email))
		})
		ginkgo.It("Returns error if the user is not logged in", func() {
			ctx := getGRPCAuthenticatedContext("INVALID_TOKEN")
			grpcUser, err := grpcAuthClient.GetUser(ctx, &pEmpty.Empty{})
			gomega.Expect(grpcUser).To(gomega.BeNil())
			gomega.Expect(err).To(gomega.HaveOccurred())
		})
	})
})
