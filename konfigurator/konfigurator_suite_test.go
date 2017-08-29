package konfigurator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestKonfigurator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Konfigurator Test Suite")
}

var _ = BeforeSuite(func() {
	httpmock.Activate()
})

var _ = BeforeEach(func() {
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

const MockResponseOidcConfiguration = `{
	"authorization_endpoint": "http://example.com/auth",
	"token_endpoint": "http://example.com/token",
	"userinfo_endpoint": "http://example.com/user-info",
	"issuer": "http://example.com",
	"jwks_uri": "http://example.com/jwks"
}`

const MockResponseToken = `{
	"access_token": "access",
	"token_type": "type",
	"refresh_token": "refresh",
	"expires_in" : "5",
	"expires": "5",
	"id_token": "super_id_token"
}`

const MockResponseMissingToken = `{
	"access_token": "access",
	"token_type": "type",
	"refresh_token": "refresh",
	"expires_in" : "5",
	"expires": "5"
}`

const (
	adfsHostUrl           = "http://example.com"
	clientID              = "fake-client-id"
	localPort             = "9999"
	localRedirectEndpoint = "/redirect-endpoint"
)
