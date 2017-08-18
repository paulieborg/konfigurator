package konfigurator_test

import (
	"bytes"
	"os"

	. "github.com/MYOB-Technology/konfigurator/konfigurator"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockReadWriteCloser struct {
	output         *bytes.Buffer
	isClosedCalled bool
	err            error
}

func (m *MockReadWriteCloser) Read(p []byte) (int, error) {
	return m.output.Read(p)
}

func (m *MockReadWriteCloser) Write(p []byte) (int, error) {
	return m.output.Write(p)
}

func (m *MockReadWriteCloser) Close() error {
	m.isClosedCalled = true
	return m.err
}

var _ = Describe("KubeConfig", func() {
	var (
		konfig *KubeConfig
		err    error
	)

	BeforeEach(func() {
		konfig, err = NewKubeConfig("123", "example.com", os.Stdout)
	})

	Context("creating a new KubeConfig", func() {
		It("should have nil error", func() {
			Expect(err).To(BeNil())
		})
		It("should have a KubeConfig populated struct", func() {
			Expect(konfig).NotTo(BeNil())
		})
		It("should have Stdout in struct", func() {
			Expect(konfig.Output).To(Equal(os.Stdout))
		})
	})

	Context("Generate config content", func() {
		BeforeEach(func() {
			konfig.Output = &MockReadWriteCloser{
				output: bytes.NewBufferString(""),
			}
			konfig.Generate("GOHAN")
		})

		It("should have the token in the output", func() {
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
			mockOutput, _ := konfig.Output.(*MockReadWriteCloser)
			Expect(mockOutput.output.String()).To(Equal(expectedContent))
		})

		It("should close output handle", func() {
			mockOutput, _ := konfig.Output.(*MockReadWriteCloser)
			Expect(mockOutput.isClosedCalled).To(BeTrue())
		})
	})
})
