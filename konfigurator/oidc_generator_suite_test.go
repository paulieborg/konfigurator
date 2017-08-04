package konfigurator_test

import (
	. "github.com/onsi/ginkgo"
	"gopkg.in/jarcoal/httpmock.v1"
)

var _ = BeforeSuite(func() {
	httpmock.Activate()
})

var _ = BeforeEach(func() {
	httpmock.Reset()
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})
