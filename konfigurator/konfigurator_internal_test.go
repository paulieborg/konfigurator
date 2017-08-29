package konfigurator

import (
	"net/http"
	"testing"
)

func TestStartHTTPServer(t *testing.T) {
	konfig := &Konfigurator{
		config: &OidcGenerator{
			localURL:              "localhost:9000",
			localRedirectEndpoint: "/some-endpoint",
		},
		tokenRetrieved: nil,
		state:          "some-state",
		kubeConfig:     &KubeConfig{},
	}

	server := konfig.startHTTPServer()
	if server == nil {
		t.Fatal("Expected server to not be nil")
	}

	res, _ := http.Head("http://localhost:9000/")
	if res.StatusCode != 302 {
		t.Fatalf("expected root to redirect, got %d", res.StatusCode)
	}

	res, _ = http.Head("http://localhost:9000/favicon.ico")
	if res.StatusCode != 204 {
		t.Fatalf("expected favicon.ico to return 204 (no content) status, got %d", res.StatusCode)
	}

	res, _ = http.Head("http://localhost:9000/some-endpoint")
	if res.StatusCode != 200 {
		t.Fatalf("expected redirect endpoint to return 200 status code, got %d", res.StatusCode)
	}
}
