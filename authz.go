package authz

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"

	"github.com/casbin/casbin/v3"
)

// CredentialValidator validates a username and password from HTTP Basic Auth.
// Assign a custom implementation before using the plugin.
// If nil, any non-empty username is accepted without password verification.
var CredentialValidator func(username, password string) bool

func init() {
	caddy.RegisterModule(Authorizer{})
	httpcaddyfile.RegisterHandlerDirective("authz", parseCaddyfile)
}

type Authorizer struct {
	AuthConfig struct {
		ModelPath  string `json:"model_path"`
		PolicyPath string `json:"policy_path"`
	} `json:"auth_config"`
	Enforcer *casbin.Enforcer `json:"-"`
}

// CaddyModule returns the Caddy module information.
func (Authorizer) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.authz",
		New: func() caddy.Module { return new(Authorizer) },
	}
}

// Provision implements caddy.Provisioner.
func (a *Authorizer) Provision(ctx caddy.Context) error {
	e, err := casbin.NewEnforcer(a.AuthConfig.ModelPath, a.AuthConfig.PolicyPath)
	if err != nil {
		return err
	}
	a.Enforcer = e
	return nil
}

// Validate implements caddy.Validator.
func (a *Authorizer) Validate() error {
	if a.Enforcer == nil {
		return fmt.Errorf("no Enforcer")
	}
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (a Authorizer) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	if CredentialValidator != nil && !CredentialValidator(username, password) {
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	allowed, err := a.CheckPermission(username, r)
	if err != nil {
		return err
	}
	if !allowed {
		w.WriteHeader(http.StatusForbidden)
		return nil
	}

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (a *Authorizer) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if !d.NextArg() {
			return d.ArgErr()
		}
		a.AuthConfig.ModelPath = d.Val()
		if !d.NextArg() {
			return d.ArgErr()
		}
		a.AuthConfig.PolicyPath = d.Val()
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Authorizer.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Authorizer
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}

// CheckPermission checks if the given user is allowed to access the resource.
func (a *Authorizer) CheckPermission(username string, r *http.Request) (bool, error) {
	return a.Enforcer.Enforce(username, r.URL.Path, r.Method)
}

// Interface guards
var (
	_ caddy.Provisioner           = (*Authorizer)(nil)
	_ caddy.Validator             = (*Authorizer)(nil)
	_ caddyhttp.MiddlewareHandler = (*Authorizer)(nil)
	_ caddyfile.Unmarshaler       = (*Authorizer)(nil)
)
