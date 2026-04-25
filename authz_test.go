package authz

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/caddyserver/caddy/v2/caddytest"
)

// testPasswords simulates a user database for tests.
var testPasswords = map[string]string{
	"alice": "alice_pass",
	"bob":   "bob_pass",
	"cathy": "cathy_pass",
}

var tester *caddytest.Tester

func testRequest(t *testing.T, user string, password string, path string, method string, code int) {
	t.Helper()
	req, err := http.NewRequest(method, fmt.Sprintf("http://localhost:9080%s", path), nil)
	if err != nil {
		t.Fatalf("unable to create request %s", err)
	}
	req.SetBasicAuth(user, password)
	tester.AssertResponse(req, code, "")
}

func initTester(t *testing.T) {
	CredentialValidator = func(username, password string) bool {
		return testPasswords[username] == password
	}

	tester = caddytest.NewTester(t)
	tester.InitServer(`
	{
		http_port     9080
		https_port    9443
		admin localhost:2999
	}
	localhost:9080 {
		route /* {
			authz "authz_model.conf" "authz_policy.csv"
			respond ""
		}
	}`, "caddyfile")
}

func TestBasic(t *testing.T) {
	initTester(t)

	testRequest(t, "alice", "alice_pass", "/dataset1/resource1", "GET", 200)
	testRequest(t, "alice", "alice_pass", "/dataset1/resource1", "POST", 200)
	testRequest(t, "alice", "alice_pass", "/dataset1/resource2", "GET", 200)
	testRequest(t, "alice", "alice_pass", "/dataset1/resource2", "POST", 403)
}

func TestWrongPassword(t *testing.T) {
	initTester(t)

	testRequest(t, "alice", "wrong_pass", "/dataset1/resource1", "GET", 401)
	testRequest(t, "bob", "wrong_pass", "/dataset2/resource1", "GET", 401)
	testRequest(t, "cathy", "wrong_pass", "/dataset1/item", "DELETE", 401)
}

func TestPathWildcard(t *testing.T) {
	initTester(t)

	testRequest(t, "bob", "bob_pass", "/dataset2/resource1", "GET", 200)
	testRequest(t, "bob", "bob_pass", "/dataset2/resource1", "POST", 200)
	testRequest(t, "bob", "bob_pass", "/dataset2/resource1", "DELETE", 200)
	testRequest(t, "bob", "bob_pass", "/dataset2/resource2", "GET", 200)
	testRequest(t, "bob", "bob_pass", "/dataset2/resource2", "POST", 403)
	testRequest(t, "bob", "bob_pass", "/dataset2/resource2", "DELETE", 403)

	testRequest(t, "bob", "bob_pass", "/dataset2/folder1/item1", "GET", 403)
	testRequest(t, "bob", "bob_pass", "/dataset2/folder1/item1", "POST", 200)
	testRequest(t, "bob", "bob_pass", "/dataset2/folder1/item1", "DELETE", 403)
	testRequest(t, "bob", "bob_pass", "/dataset2/folder1/item2", "GET", 403)
	testRequest(t, "bob", "bob_pass", "/dataset2/folder1/item2", "POST", 200)
	testRequest(t, "bob", "bob_pass", "/dataset2/folder1/item2", "DELETE", 403)
}

func TestRBAC(t *testing.T) {
	initTester(t)

	// cathy can access all /dataset1/* resources via all methods because it has the dataset1_admin role.
	testRequest(t, "cathy", "cathy_pass", "/dataset1/item", "GET", 200)
	testRequest(t, "cathy", "cathy_pass", "/dataset1/item", "POST", 200)
	testRequest(t, "cathy", "cathy_pass", "/dataset1/item", "DELETE", 200)
	testRequest(t, "cathy", "cathy_pass", "/dataset2/item", "GET", 403)
	testRequest(t, "cathy", "cathy_pass", "/dataset2/item", "POST", 403)
	testRequest(t, "cathy", "cathy_pass", "/dataset2/item", "DELETE", 403)
}
