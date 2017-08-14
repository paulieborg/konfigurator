package konfigurator

import (
	"context"
	"errors"

	oidc "github.com/coreos/go-oidc"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

// Struct that deals with OIDC information such as the redirect endpoint and all the Oauth2 config
type OidcGenerator struct {
	ctx                   context.Context
	config                oauth2.Config
	localUrl              string
	localRedirectEndpoint string
}

// Contructor for OidcGenerator which uses a default background context and 'localhost' for the redirectUrl
func NewOidcGenerator(adfsHostUrl, clientId, localPort, localRedirectEndpoint string) (*OidcGenerator, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, adfsHostUrl)
	if err != nil {
		return nil, err
	}

	localUrl := "localhost:" + localPort
	return &OidcGenerator{
		ctx: ctx,
		config: oauth2.Config{
			ClientID:    clientId,
			RedirectURL: "http://" + localUrl + localRedirectEndpoint,
			Endpoint:    provider.Endpoint(),
		},
		localUrl:              localUrl,
		localRedirectEndpoint: localRedirectEndpoint,
	}, nil
}

// Simply allows the same method call to be passed on to the underlying Oauth2 config struct
func (o *OidcGenerator) AuthCodeURL(state string) string {
	return o.config.AuthCodeURL(state)
}

func (o *OidcGenerator) openBrowser() {
	open.Run("http://" + o.localUrl)
}

// Retrieves the Oauth2 token from the request and extracts the "id_token" part of it
func (o *OidcGenerator) GetToken(code string) (string, error) {
	oauth2Token, err := o.config.Exchange(o.ctx, code)
	if err != nil {
		return "", err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", errors.New("missing id_token from oauth2 token")
	}

	return rawIDToken, nil
}
