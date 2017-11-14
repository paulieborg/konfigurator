/*
Package konfigurator provides a library for generating Kubernetes config files by means of OpenID connect authentication.
It will do an OIDC Token exchange to the Host given and create a configuration file with all the information provided as
well as the token retrieved. By default, konfigurator will output the contents of the file to `stdout`, this can be overridden
with the `-o|--output` flag.

NOTE: that this tool will start a local webserver in the provided port to be able to handle the callback from the OpenID Connect
protocol, so it is important to make sure the port provided is not in use by the host.
*/
package konfigurator

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

type Konfigurator struct {
	config         *OidcGenerator
	tokenRetrieved chan int
	state          string
	nonceValue     string
	kubeConfig     *KubeConfig
}

// NewKonfigurator creates a file and a uuid to use as a state to check MITM attacks and returns a new Konfigurator struct.
func NewKonfigurator(oidcHost, oidcClientID, oidcClientPort, oidcClientRedirectEndpoint, kubeCa, kubeAPIURL, kubeNamespace, outputFilePath string) (*Konfigurator, error) {
	config, err := NewOidcGenerator(oidcHost, oidcClientID, oidcClientPort, oidcClientRedirectEndpoint)
	if err != nil {
		return nil, err
	}

	uid, _ := uuid.NewV4()
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

	kubeConfig, err := NewKubeConfig(kubeCa, kubeAPIURL, kubeNamespace, fileHandle)
	if err != nil {
		return nil, err
	}

	return &Konfigurator{
		config,
		make(chan int, 1),
		uid.String(),
		string(rand.New(rand.NewSource(time.Now().UnixNano())).Int()),
		kubeConfig,
	}, nil
}

// Orchestrate will start a local web server based on parameters from the constructor,
// will open a browser and initiate the authentication process. Once the process is done,
// it will output the kubernetes config file to the output file path (or stdout of that is empty)
// and close the web server. The webserver will only be closed once the authentication succeeds.
func (k *Konfigurator) Orchestrate() error {
	server := k.startHTTPServer()
	k.config.OpenBrowser()

	// block
	<-k.tokenRetrieved

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		return err // failure/timeout shutting down the server gracefully
	}

	return nil
}

func (k *Konfigurator) startHTTPServer() *http.Server {
	srv := &http.Server{Addr: k.config.localURL}

	http.HandleFunc("/", k.rootHandler)
	http.HandleFunc("/favicon.ico", k.noContentHandler)
	http.HandleFunc(k.config.localRedirectEndpoint, tokenRedirectCallbackHandler)
	http.HandleFunc("/auth/js/redirect", k.callbackHandler)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
		return
	}()

	return srv
}

func (k *Konfigurator) rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, k.config.AuthCodeURL(k.state, k.nonceValue), http.StatusFound)
}

func (k *Konfigurator) noContentHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (k *Konfigurator) callbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != k.state {
		log.Printf("URL State did not match: expected %s, got %s", k.state, r.URL.Query().Get("state"))
		return
	}

	token := r.URL.Query().Get("id_token")

	k.kubeConfig.Generate(token)
	io.WriteString(w, httpContent)
	// unblock
	k.tokenRetrieved <- 1
	return
}

func tokenRedirectCallbackHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, httpTokenRedirectContent)
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

var httpTokenRedirectContent = `
<html>
	<script>
		(function() {
			var hash = location.hash.slice(1);
			window.location = "/auth/js/redirect?" + hash;
		})()
	</script>
</html>
`
