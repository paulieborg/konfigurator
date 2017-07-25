package konfigurator

import (
	"context"
	"errors"

	oidc "github.com/coreos/go-oidc"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

type oidcGenerator struct {
	ctx                   context.Context
	config                oauth2.Config
	localUrl              string
	localRedirectEndpoint string
}

func NewOidcGenerator(adfsHostUrl, clientId, localPort, localRedirectEndpoint string) (*oidcGenerator, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, adfsHostUrl)
	if err != nil {
		return nil, err
	}

	localUrl := "localhost:" + localPort
	return &oidcGenerator{
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

func (o *oidcGenerator) AuthCodeURL(state string) string {
	return o.config.AuthCodeURL(state)
}

func (o *oidcGenerator) openBrowser() {
	open.Run("http://" + o.localUrl)
}

func (o *oidcGenerator) GetToken(code string) (string, error) {
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
