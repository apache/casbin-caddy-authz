# casbin-caddy-authz

[![Go](https://github.com/apache/casbin-caddy-authz/actions/workflows/ci.yml/badge.svg)](https://github.com/apache/casbin-caddy-authz/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/apache/casbin-caddy-authz/badge.svg?branch=master)](https://coveralls.io/github/apache/casbin-caddy-authz?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/apache/casbin-caddy-authz)](https://goreportcard.com/report/github.com/apache/casbin-caddy-authz)
[![Godoc](https://godoc.org/github.com/apache/casbin-caddy-authz?status.svg)](https://godoc.org/github.com/apache/casbin-caddy-authz)

casbin-caddy-authz is an authorization middleware for [Caddy](https://github.com/caddyserver/caddy), based on [Apache Casbin](https://github.com/apache/casbin). It controls access to your web resources by enforcing authorization policies defined with Apache Casbin.

## Installation

```
go get github.com/apache/casbin-caddy-authz/v2
```

## Simple Example

```go
package main

import (
    "github.com/caddyserver/caddy/v2"
    _ "github.com/apache/casbin-caddy-authz/v2"
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

1. Path to the Apache Casbin **model file** (`.conf`) — describes the access control model (ACL, RBAC, ABAC, etc.)
2. Path to the Apache Casbin **policy file** (`.csv`) — describes the authorization rules

For how to write these files, refer to the [Apache Casbin documentation](https://casbin.apache.org/docs/get-started).

## Security: Authentication vs Authorization

> **⚠️ Important:** This plugin handles **authorization only** — it does NOT validate passwords or verify user identity.
>
> You **must** place an authentication middleware **before** this plugin to verify credentials. Without it, anyone can set an arbitrary `Authorization` header and impersonate any user.

This plugin is designed to be used alongside a dedicated authentication layer:

```
localhost:8080 {
    route {
        basicauth {                                      # ← Step 1: Authentication (validates credentials)
            alice $2a$14$Zkx19XLiW6VYouLHR5NmfOFU0z2GTNmpkT/5qqR7hx4IjWJPDhjvG
        }
        authz "authz_model.conf" "authz_policy.csv"    # ← Step 2: Authorization (checks permissions)
        respond "Hello, world!"
    }
}
```

Caddy's built-in [`basicauth`](https://caddyserver.com/docs/caddyfile/directives/basic_auth) directive is recommended for HTTP Basic Authentication. For other schemes (OAuth, JWT, LDAP, etc.), use the appropriate authentication plugin and ensure it runs before `authz`.

If you need to validate credentials programmatically, set the `CredentialValidator` hook before the server starts:

```go
import authz "github.com/apache/casbin-caddy-authz/v3"

func init() {
    authz.CredentialValidator = func(username, password string) bool {
        // query your database or LDAP here
        return myDB.CheckPassword(username, password)
    }
}
```

## How Access Control Works

Authorization is determined based on `{subject, object, action}`:

| Field | Meaning |
|-------|---------|
| `subject` | The authenticated user name (extracted from HTTP Basic Auth header) |
| `object` | The URL path of the requested resource, e.g. `dataset1/item1` |
| `action` | The HTTP method, e.g. `GET`, `POST`, `PUT`, `DELETE` |

## Working Example

1. Build Caddy with this plugin using [xcaddy](https://github.com/caddyserver/xcaddy):

    ```bash
    xcaddy build --with github.com/apache/casbin-caddy-authz/v2
    ```

2. Place your Apache Casbin model file [authz_model.conf](https://github.com/apache/casbin-caddy-authz/blob/master/authz_model.conf) and policy file [authz_policy.csv](https://github.com/apache/casbin-caddy-authz/blob/master/authz_policy.csv) in a known directory.

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

- [Apache Casbin](https://github.com/apache/casbin)

## License

This project is under Apache 2.0 License. See the [LICENSE](LICENSE) file for the full license text.
