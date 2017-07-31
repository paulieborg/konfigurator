package konfigurator_test

import (
	"bytes"
	"os"
	"testing"

	. "github.com/MYOB-Technology/konfigurator/konfigurator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKonfigurator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Konfigurator Suite")
}

var _ = Describe("KubeConfig", func() {
	var (
		konfig *KubeConfig
		err    error
	)

	BeforeEach(func() {
		konfig, err = NewKubeConfig("123", "example.com")
	})

	Context("creating a new KubeConfig", func() {
		It("should have nil error", func() {
			Expect(err).To(BeNil())
		})
		It("should have a KubeConfig populated struct", func() {
			Expect(konfig).NotTo(BeNil())
		})
		It("should have Stdout in struct", func() {
			Expect(konfig.File).To(Equal(os.Stdout))
		})
	})

	Context("Generate config content", func() {
		It("should have the token in the output", func() {
			konfig.File = bytes.NewBufferString("")
			konfig.Generate("GOHAN")
			expectedContent := `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: 123
    server: https://api.example.com
  name: example.com
contexts:
- context:
    cluster: example.com
    user: OIDCUser
  name: example.com
current-context: example.com
kind: Config
preferences: {}
users:
- name: OIDCUser
  user:
    token: GOHAN
`
			buf, _ := konfig.File.(*bytes.Buffer)
			Expect(buf.String()).To(Equal(expectedContent))
		})
	})
})
