package integration

import (
	"net/http"

	"errors"
	"github.com/axiomzen/authentication/models"
	"github.com/axiomzen/authentication/routes"
	"github.com/axiomzen/compare"
	"github.com/axiomzen/yawgh"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

//

var _ = ginkgo.Describe("Ping", func() {
	ginkgo.Context("Pingdom", func() {
		ginkgo.It("should return status okay without version in path and no api token", func() {
			// get response
			var pingBack models.Ping
			statusCode, err := yawgh.New().
				Transport("http").
				DomainHost(theConf.TestDomainHost).
				Port(uint(theConf.Port)).
				Marshaler(marshaler).
				Unmarshaler(unmarshaler).
				ResponseInterceptor(locationChecker).
				ResponseBody(&pingBack).
				Get(routes.ResourcePing).
				Do()

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			// check status code
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			// check result
			expectedPing := models.Ping{Ping: "pong"}
			gomega.Ω(compare.New().DeepEquals(pingBack, expectedPing, "Ping")).Should(gomega.Succeed())
		})

		ginkgo.It("should return status okay without version in path and api token", func() {
			// get response
			var pingBack models.Ping

			statusCode, err := yawgh.New().
				Transport("http").
				DomainHost(theConf.TestDomainHost).
				Port(uint(theConf.Port)).
				Marshaler(marshaler).
				Unmarshaler(unmarshaler).
				ResponseInterceptor(locationChecker).
				ResponseBody(&pingBack).
				Header(theConf.APITokenHeader, theConf.APIToken).
				Get(routes.ResourcePing).
				Do()

			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			// check status code
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			// check result
			expectedPing := models.Ping{Ping: "pong"}
			gomega.Ω(compare.New().DeepEquals(pingBack, expectedPing, "Ping")).Should(gomega.Succeed())
		})

		// todo: I think we want it with api token only
		ginkgo.It("should return status okay with version in path", func() {
			// get response
			var pingBack models.Ping
			statusCode, err := TestRequestV1().Get(routes.ResourcePing).ResponseBody(&pingBack).Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			// check status code
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			// check result
			expectedPing := models.Ping{Ping: "pong"}
			gomega.Ω(compare.New().DeepEquals(pingBack, expectedPing, "Ping")).Should(gomega.Succeed())
		})

		ginkgo.It("should return status ok with no content-type", func() {
			var fn responseIntFunc = func(r *http.Response, body []byte, contentType string) error {
				if r.Header.Get("Content-Type") != "text/html" {
					return errors.New("Unexpected content type!")
				}
				return nil
			}
			statusCode, err := TestRequestV1().
				ContentType("").
				Header("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
				Get(routes.ResourcePing).
				//RequestInterceptor(printRequest).
				ResponseInterceptor(fn).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			// check status code
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			// (optional) check result
			// expectedPing := "&models.Ping{Ping:\"pong\"}"
			// gomega.Ω(compare.New().DeepEquals(pingBack, expectedPing, "Ping")).Should(gomega.Succeed())
		})

		ginkgo.It("should return status ok with custom content-type", func() {
			var pingBack models.Ping
			statusCode, err := TestRequestV1().
				ContentType("application/json; charset=utf-8").
				//Header("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8").
				Get(routes.ResourcePing).
				RequestInterceptor(printRequest).
				ResponseBody(&pingBack).
				Do()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			// check status code
			gomega.Expect(statusCode).To(gomega.Equal(http.StatusOK))
			// check result
			expectedPing := models.Ping{Ping: "pong"}
			gomega.Ω(compare.New().DeepEquals(pingBack, expectedPing, "Ping")).Should(gomega.Succeed())
		})
	})
})
