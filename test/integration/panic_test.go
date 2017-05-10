package integration

import (
	"net/http"

	"github.com/axiomzen/authentication/constants"
	"github.com/axiomzen/authentication/models"
	"github.com/axiomzen/authentication/routes"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

// Panic testing routes
var _ = ginkgo.Describe("Panic", func() {
	ginkgo.Context("Panic Tests", func() {
		ginkgo.It("should return status 500 error upon panic", func() {
			var errorResp models.ErrorResponse
			statusCode, err := TestRequestV1().
				Get(routes.ResourceTest + routes.ResourcePanic).
				ErrorResponseBody(&errorResp).
				Do()
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusInternalServerError))
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(errorResp.Code).To(gomega.Equal(constants.APIPanic))

		})
	})
})
