package konfigurator_test

import (
	"net/http"
	"os"

	. "github.com/MYOB-Technology/konfigurator/konfigurator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var _ = Describe("Konfigurator", func() {
	var (
		konfig                     *Konfigurator
		err                        error
		oidcGeneratorMockResponder func()
	)

	oidcGeneratorMockResponder = func() {
		httpmock.RegisterResponder(
			"GET",
			adfsHostUrl+"/.well-known/openid-configuration",
			httpmock.NewStringResponder(http.StatusOK, MockResponseOidcConfiguration),
		)
	}

	Describe("A valid Konfigurator", func() {
		BeforeEach(func() {
			oidcGeneratorMockResponder()
			konfig, err = NewKonfigurator(adfsHostUrl, clientID, localPort, localRedirectEndpoint, "CA", "api.url.com", "", "/tmp/path")
		})

		Context("creating a new Konfigurator", func() {
			It("should have nil error", func() {
				Expect(err).To(BeNil())
			})

			It("should have a randonly generated state", func() {
				Expect(konfig).NotTo(BeNil())
			})
		})
	})

	Describe("An invalid Konfigurator", func() {
		Context("Error creating an OidcGenerator", func() {
			It("should return an error", func() {
				konfig, err = NewKonfigurator(adfsHostUrl, clientID, localPort, localRedirectEndpoint, "CA", "api.url.com", "", "/tmp/path")
				Expect(err).NotTo(BeNil())
				Expect(konfig).To(BeNil())
			})
		})

		Context("Error creating a new file", func() {
			It("should error when creating a new file", func() {
				oidcGeneratorMockResponder()
				konfig, err = NewKonfigurator(adfsHostUrl, clientID, localPort, localRedirectEndpoint, "cd", "api.url.com", "someNamespace", "/tmp/somepath/!~!~!@#!#@that/wont/exists/config")
				_, ok := err.(*os.PathError)
				Expect(ok).To(BeTrue())
				Expect(konfig).To(BeNil())
			})
		})
	})
})
