# Kubectl OIDC Config Generator

[![Build status](https://badge.buildkite.com/10540b85a34e726f839daf35543aec5e484dba8a32a63f3491.svg)](https://buildkite.com/myob/konfigurator)

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

