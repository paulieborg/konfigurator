package konfigurator

import (
	"html/template"
	"io"
	"os"
)

type KubeConfig struct {
	CA   string
	URL  string
	tmpl *template.Template
	File io.ReadWriter
}

var content = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{.CA}}
    server: https://api.{{.URL}}
  name: {{.URL}}
contexts:
- context:
    cluster: {{.URL}}
    user: OIDCUser
  name: {{.URL}}
current-context: {{.URL}}
kind: Config
preferences: {}
users:
- name: OIDCUser
  user:
    token: {{.Token}}
`

type configData struct {
	CA    string
	URL   string
	Token string
}

func NewKubeConfig(ca, url string) (*KubeConfig, error) {
	tmpl, err := template.New("config").Parse(content)
	if err != nil {
		return nil, err
	}

	return &KubeConfig{
		ca,
		url,
		tmpl,
		os.Stdout,
	}, nil
}

func (k *KubeConfig) Generate(token string) error {
	err := k.tmpl.Execute(k.File, configData{
		k.CA,
		k.URL,
		token,
	})
	if err != nil {
		return err
	}

	return nil
}
