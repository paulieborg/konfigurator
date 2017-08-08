package konfigurator

import (
	"log"
	"net/http"
	"os"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type konfigurator struct {
	config         *OidcGenerator
	tokenRetrieved bool
	state          string
	kubeConfig     *KubeConfig
}

func NewKonfigurator(oidcHost, oidcClientId, oidcClientPort, oidcClientRedirectEndpoint, kubeCa, kubeApiUrl, outputFilePath string) (*konfigurator, error) {
	config, err := NewOidcGenerator(oidcHost, oidcClientId, oidcClientPort, oidcClientRedirectEndpoint)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	fileHandle := os.Stdout
	if outputFilePath != "" {
		fileHandle, err = os.Create(outputFilePath)
		if err != nil {
			return nil, err
		}
	}

	kubeConfig, err := NewKubeConfig(kubeCa, kubeApiUrl, fileHandle)
	if err != nil {
		return nil, err
	}

	return &konfigurator{
		config,
		false,
		uid.String(),
		kubeConfig,
	}, nil
}

func (k *konfigurator) Orchestrate() error {
	server := k.startHttpServer()
	k.config.openBrowser()

	for !k.tokenRetrieved {
		time.Sleep(1 * time.Second)
	}

	err := server.Shutdown(nil)
	if err != nil {
		return err // failure/timeout shutting down the server gracefully
	}

	return nil
}

func (k *konfigurator) startHttpServer() *http.Server {
	srv := &http.Server{Addr: k.config.localUrl}

	http.HandleFunc("/", k.rootHandler)
	http.HandleFunc("/favicon.ico", k.noContentHandler)
	http.HandleFunc(k.config.localRedirectEndpoint, k.callbackHandler)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	return srv
}

func (k *konfigurator) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, k.config.AuthCodeURL(k.state), http.StatusFound)
}

func (k *konfigurator) noContentHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (k *konfigurator) callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != k.state {
		panic("state did not match")
	}

	token, err := k.config.GetToken(r.URL.Query().Get("code"))
	if err != nil {
		panic(err)
	}

	k.kubeConfig.Generate(token)
	w.Write([]byte(httpContent))
	k.tokenRetrieved = true
}

var httpContent = `
<html>
    <body>
        Token retrieved successfully.
        This tab will close soon.

        <script>
            setTimeout(function() {
                window.close();
            }, 2000);
        </script>
    </body>
</html>
`
