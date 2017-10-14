# Kubectl OIDC Config Generator

[![Build Status](https://travis-ci.org/MYOB-Technology/konfigurator.svg?branch=master)](https://travis-ci.org/MYOB-Technology/konfigurator)
[![Coverage Status](https://coveralls.io/repos/github/MYOB-Technology/konfigurator/badge.svg?branch=master)](https://coveralls.io/github/MYOB-Technology/konfigurator?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/MYOB-Technology/Konfigurator)](https://goreportcard.com/report/github.com/MYOB-Technology/Konfigurator)

This tool generates a `kubeconfig` file with the OpenID Connect token authentication, allowing a user to use `kubectl` with OIDC configured clusters.

## Usage

```bash
❯ ./build/konfigurator -h
NAME:
   konfigurator - generate a kubeconfig file with OIDC token

USAGE:
   konfigurator [global options] command [command options] [arguments...]

VERSION:
   SNAPSHOT

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --client-id value, -c value              oidc provider's client id
   --host value, -u value                   oidc provider's host url
   --port value, -p value                   port for the local server (default: "9000")
   --redirect-endpoint value, -r value      endpoint for oidc redirect (default: "/oauth2/callback")
   --api value, -a value                    url for kubernetes api
   --certificate-authority value, -s value  certificate authority cert for kubernetes api
   --help, -h                               show help
   --version, -v                            print the version
```

### Authenticating against an OAuth2 Auth0 setup

```bash
❯ konfigurator --client-id SomeClientIDHere --host https://gugahoi.au.auth0.com/ --api your.k8s.domain.com -s SOME-CA-CONTENT-HERE

apiVersion: v1
clusters:
- cluster:
    certificate-authority-data:
    SOME-CA-CONTENT-HERE
    server: https://api.your.k8s.domain.com
  name: your.k8s.domain.com
contexts:
- context:
    cluster: your.k8s.domain.com
    namespace: default
    user: OIDCUser
  name: your.k8s.domain.com
current-context: your.k8s.domain.com
kind: Config
preferences: {}
users:
- name: OIDCUser
  user:
    token: eyJ0eXWithGeneratedContentF1MtJIH4Vg
```

### Using it as a Library

```go
// main.go
package main

import gihtub.com/MYOB-Technology/konfigurator/konfigurator

var (
    clientID         = "123456-9999-9876-6789-123456789"
    host             = "https://gugahoi.au.auth0.com/"
    port             = "9000"
    redirectEndpoint = "/oauth2/callback"
    cluster = konfigurator.KubeConfig{
        CA:  "SOME-CA-CONTENT-HERE",
        URL: "your.k8s.domain.com",
    }
)

function main(){
    k, err := konfigurator.NewKonfigurator(host, clientID, port, redirectEndpoint, cluster.CA, cluster.URL, "default", "~/.kube/adfs-config")
    if nil != err {
        return 1, err
    }
    err = k.Orchestrate()
    if err != nil {
        return 1, err // failure/timeout shutting down the server gracefully
    }
    return 0, nil
}
```
