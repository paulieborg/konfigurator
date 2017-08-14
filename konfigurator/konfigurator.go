package konfigurator

import (
	"io"
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

// NewKonfigurator creates a file and a uuid to use as a state to check MITM attacks and returns a new Konfigurator struct.
func NewKonfigurator(oidcHost, oidcClientID, oidcClientPort, oidcClientRedirectEndpoint, kubeCa, kubeAPIURL, outputFilePath string) (*konfigurator, error) {
	config, err := NewOidcGenerator(oidcHost, oidcClientID, oidcClientPort, oidcClientRedirectEndpoint)
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

	kubeConfig, err := NewKubeConfig(kubeCa, kubeAPIURL, fileHandle)
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
	server := k.startHTTPServer()
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

func (k *konfigurator) startHTTPServer() *http.Server {
	srv := &http.Server{Addr: k.config.localURL}

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
		log.Printf("URL State did not match: expected %s, got %s", k.state, r.URL.Query().Get("state"))
		return
	}

	token, err := k.config.GetToken(r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("Failed extracting token: %s", err)
		return
	}

	k.kubeConfig.Generate(token)
	io.WriteString(w, httpContent)
	k.tokenRetrieved = true
	return
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
