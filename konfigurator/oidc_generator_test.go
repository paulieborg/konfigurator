package konfigurator_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	. "github.com/MYOB-Technology/konfigurator/konfigurator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/jarcoal/httpmock.v1"
)

var _ = Describe("OidcGenerator", func() {
	var (
		konfig *OidcGenerator
		err    error
	)

	Describe("valid OidcGenerator", func() {
		BeforeEach(func() {
			httpmock.RegisterResponder(
				"GET",
				hostUrl+"/.well-known/openid-configuration",
				httpmock.NewStringResponder(http.StatusOK, MockResponseOidcConfiguration),
			)

			konfig, err = NewOidcGenerator(hostUrl, clientID, localPort, localRedirectEndpoint)
		})

		Context("creating a new OidcGenerator", func() {
			It("should have nil error", func() {
				Expect(err).To(BeNil())
			})

			It("should have a OidcGenerator populated struct", func() {
				Expect(konfig).NotTo(BeNil())
			})
		})

		Context("AuthCodeURL", func() {
			It("should return a url with the given state, nonce and scope", func() {
				state := "some-random-state"
				nonceSecretValue := "cryptography-yay"
				Expect(konfig.AuthCodeURL(state, nonceSecretValue)).To(
					Equal(fmt.Sprintf(
						"http://example.com/auth?client_id=%s&nonce=%s&redirect_uri=%s&response_type=id_token&scope=openid&state=%s",
						clientID,
						nonceSecretValue,
						url.QueryEscape("http://localhost:"+localPort+localRedirectEndpoint),
						state,
					)),
				)
			})
		})

		Context("GetToken", func() {
			mockToken := func(mockResponse string) {
				httpmock.RegisterResponder(
					"POST",
					"http://example.com/token",
					httpmock.NewStringResponder(http.StatusOK, mockResponse),
				)
			}

			It("should get a JWT token", func() {
				mockToken(MockResponseToken)
				token, err := konfig.GetToken("vegeta-is-a-sayan-prince")

				Expect(err).To(BeNil())
				Expect(token).To(Equal("super_id_token"))
			})

			It("should fail to get id_token from oauthToken", func() {
				mockToken(MockResponseMissingToken)
				token, err := konfig.GetToken("vegeta-is-a-sayan-prince")

				Expect(err).NotTo(BeNil())
				Expect(token).To(Equal(""))
			})

			It("should fail to call Exchange", func() {
				token, err := konfig.GetToken("failure-mode-token")
				Expect(err).NotTo(BeNil())
				Expect(token).To(Equal(""))
			})
		})
	})

	Describe("invalid OidcGenerator", func() {
		Context("oidcGenerator Fail Scenarios", func() {
			It("should return an error when creating OidcGenerator", func() {
				konfig, err = NewOidcGenerator("some-invalid-host.com", "123", "999", "/endpoint")
				Expect(err).NotTo(BeNil())
				Expect(konfig).To(BeNil())
			})
		})
	})

	Describe("openBrowser", func() {
		Context("calling Run", func() {
			It("should call the Runner", func() {
				var err error
				var isCalled bool
				out := bytes.NewBuffer([]byte{})
				konfig := OidcGenerator{
					Run: func(string) error {
						isCalled = true
						return nil
					},
					Stream: out,
				}
				konfig.OpenBrowser()

				Expect(isCalled).To(BeTrue())
				Expect(err).To(BeNil())
				Expect(out.String()).To(BeEmpty())
			})

			It("should call a failed runner", func() {
				var err error
				var isCalled bool
				out := bytes.NewBuffer([]byte{})
				konfig := OidcGenerator{
					Run: func(string) error {
						isCalled = true
						err = fmt.Errorf("some error")
						return err
					},
					Stream: out,
				}

				konfig.OpenBrowser()

				Expect(isCalled).To(BeTrue())
				Expect(err).ToNot(BeNil())
				Expect(out).To(MatchRegexp("Go to the following url to authenticate: http://.*."))
			})
		})
	})
})
