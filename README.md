# Caddy-authz

[![Go](https://github.com/casbin/caddy-authz/actions/workflows/ci.yml/badge.svg)](https://github.com/casbin/caddy-authz/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/casbin/caddy-authz/badge.svg?branch=master)](https://coveralls.io/github/casbin/caddy-authz?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/casbin/caddy-authz)](https://goreportcard.com/report/github.com/casbin/caddy-authz)
[![Godoc](https://godoc.org/github.com/casbin/caddy-authz?status.svg)](https://godoc.org/github.com/casbin/caddy-authz)

Caddy-authz is an authorization middleware for [Caddy](https://github.com/caddyserver/caddy), based on [Casbin](https://github.com/casbin/casbin). It controls access to your web resources by enforcing authorization policies defined with Casbin.

## Installation

```
go get github.com/casbin/caddy-authz/v2
```

## Simple Example

```go
package main

import (
    "github.com/caddyserver/caddy/v2"
    _ "github.com/casbin/caddy-authz/v2"
)

func main() {
    caddy.Run(&caddy.Config{})
}
```

## Caddyfile Syntax

```
localhost {
    route {
        authz "/path/to/authz_model.conf" "/path/to/authz_policy.csv"
    }
    respond "Hello, world!"
}
```

Or using global options to control directive ordering:

```
{
    order authz before respond
}

localhost {
    authz "/path/to/authz_model.conf" "/path/to/authz_policy.csv"
    respond "Hello, world!"
}
```

The `authz` directive takes two arguments:

1. Path to the Casbin **model file** (`.conf`) — describes the access control model (ACL, RBAC, ABAC, etc.)
2. Path to the Casbin **policy file** (`.csv`) — describes the authorization rules

For how to write these files, refer to the [Casbin documentation](https://casbin.org/docs/get-started).

## How Access Control Works

Authorization is determined based on `{subject, object, action}`:

| Field | Meaning |
|-------|---------|
| `subject` | The logged-in user name (from HTTP Basic Auth header) |
| `object` | The URL path of the requested resource, e.g. `dataset1/item1` |
| `action` | The HTTP method, e.g. `GET`, `POST`, `PUT`, `DELETE` |

> **Note:** This plugin reads the user name from the HTTP `Authorization` header using Basic Auth. If you use other authentication methods (OAuth, LDAP, JWT, etc.), you will need to customize the plugin.

## Working Example

1. Build Caddy with this plugin using [xcaddy](https://github.com/caddyserver/xcaddy):

    ```bash
    xcaddy build --with github.com/casbin/caddy-authz/v2
    ```

2. Place your Casbin model file [authz_model.conf](https://github.com/casbin/caddy-authz/blob/master/authz_model.conf) and policy file [authz_policy.csv](https://github.com/casbin/caddy-authz/blob/master/authz_policy.csv) in a known directory.

3. Add the `authz` directive to your `Caddyfile`:

    ```
    localhost:8080 {
        route {
            authz "authz_model.conf" "authz_policy.csv"
        }
        respond "Hello, world!"
    }
    ```

4. Run `caddy` and enjoy.

## Getting Help

- [Casbin](https://github.com/casbin/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
