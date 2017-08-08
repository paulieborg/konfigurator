# Kubectl OIDC Config Generator

[![Build Status](https://travis-ci.org/MYOB-Technology/konfigurator.svg?branch=master)](https://travis-ci.org/MYOB-Technology/konfigurator)
[![Coverage Status](https://coveralls.io/repos/github/MYOB-Technology/konfigurator/badge.svg?branch=master)](https://coveralls.io/github/MYOB-Technology/konfigurator?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/MYOB-Technology/Konfigurator)](https://goreportcard.com/report/github.com/MYOB-Technology/Konfigurator)

This tool generates a `kubeconfig` file with the OpenID Connect token authentication, allowing a user to use `kubectl` with OIDC configured clusters.

# Usage

```
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

