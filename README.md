<!-- markdownlint-disable-next-line first-line-h1 -->
[![OpenSSFScorecard](https://api.securityscorecards.dev/projects/github.com/nginx/nginx-plus-go-client/badge)](https://scorecard.dev/viewer/?uri=github.com/nginx/nginx-plus-go-client)
[![Continuous Integration](https://github.com/nginx/nginx-plus-go-client/workflows/Continuous%20Integration/badge.svg)](https://github.com/nginx/nginx-plus-go-client/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/nginx/nginx-plus-go-client)](https://goreportcard.com/report/github.com/nginx/nginx-plus-go-client)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/nginx/nginx-plus-go-client?logo=github&sort=semver)](https://github.com/nginx/nginx-plus-go-client/releases/latest)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/nginx/nginx-plus-go-client?logo=go)
[![Slack](https://img.shields.io/badge/slack-nginxcommunity-green?logo=slack)](https://nginxcommunity.slack.com)
[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)
[![Community Support](https://badgen.net/badge/support/community/cyan?icon=awesome)](https://github.com/nginx/nginx-plus-go-client/blob/main/SECURITY.md)

# NGINX Plus Go Client

This project includes a client library for working with NGINX Plus API.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
## Table of Contents

- [About the Client](#about-the-client)
- [Compatibility](#compatibility)
- [Using the Client](#using-the-client)
- [Testing](#testing)
  - [Unit tests](#unit-tests)
  - [Integration tests](#integration-tests)
- [Contacts](#contacts)
- [Contributing](#contributing)
- [Support](#support)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## About the Client

`client/nginx.go` includes functions and data structures for working with NGINX Plus API as well as some helper
functions.

## Compatibility

This Client works against versions 4 to 9 of the NGINX Plus API. The table below shows the version of NGINX Plus where
the API was first introduced.

| API version | NGINX Plus version |
| ----------- | ------------------ |
| 4           | R18                |
| 5           | R19                |
| 6           | R20                |
| 7           | R25                |
| 8           | R27                |
| 9           | R30                |

## Using the Client

1. Import `github.com/nginx/nginx-plus-go-client/client` into your go project.
2. Use your favorite vendor tool to add this to your `/vendor` directory in your project.

## Testing

### Unit tests

```console
cd client
go test
```

### Integration tests

Prerequisites:

- Docker
- Docker Compose
- golang
- Make
- NGINX Plus license - put `nginx-repo.crt` and `nginx-repo.key` into the `docker` folder.

Run Tests:

```console
make test
```

This will build and run two NGINX Plus containers and create one docker network of type bridge, execute the client tests
against both NGINX Plus APIs, and then clean up. If it fails and you want to clean up (i.e. stop the running containers
and remove the docker network), please use `make clean`

## Contacts

We’d like to hear your feedback! If you have any suggestions or experience issues with the NGINX Plus Go Client, please
create an issue or send a pull request on GitHub. You can contact us directly via <integrations@nginx.com> or on the
[NGINX Community Slack](https://nginxcommunity.slack.com).

## Contributing

If you'd like to contribute to the project, please read our [Contributing guide](CONTRIBUTING.md).

## Support

This project is not covered by the NGINX Plus support contract.
