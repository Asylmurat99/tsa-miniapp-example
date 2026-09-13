package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	testAppID  = "34d5810c-7c36-11eb-82e7-f2189812cd57"
	testSecret = "my_secret_key"
	testUser   = `{"id":"3f1a9c40-77ad-4e0f-9c3d-11c2a0f5e881"}`
)

// newTestServer returns the api too, so a test can move the clock.
func newTestServer(t *testing.T) (*httptest.Server, *api) {
	t.Helper()
	a := &api{
		verifier: Verifier{AppID: testAppID, Secret: testSecret, Window: 300 * time.Second, Now: time.Now},
		sessions: newSessionStore(),
		log:      log.New(io.Discard, "", 0),
	}
	mux := http.NewServeMux()
	a.routes(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, a
}

func postSession(t *testing.T, srv *httptest.Server, initData string) (*http.Response, map[string]any) {
	t.Helper()
	resp, err := http.Post(srv.URL+"/api/session", "text/plain", strings.NewReader(initData))
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&body)
	resp.Body.Close()
	return resp, body
}

func TestSessionCustomer(t *testing.T) {
	srv, _ := newTestServer(t)
	initData := buildInitData(testAppID, testSecret, Fields{"auth": "customer", "user": testUser, "scope": "phone:read"})
	resp, body := postSession(t, srv, initData)
	if resp.StatusCode != 200 {
		t.Fatalf("status %d body %v", resp.StatusCode, body)
	}
	if body["auth"] != "customer" || body["user_id"] != "3f1a9c40-77ad-4e0f-9c3d-11c2a0f5e881" {
		t.Fatalf("body %v", body)
	}
	if resp.Header.Get("Set-Cookie") == "" || !strings.Contains(resp.Header.Get("Set-Cookie"), "HttpOnly") {
		t.Fatalf("cookie %q", resp.Header.Get("Set-Cookie"))
	}
}

func TestSessionGuest(t *testing.T) {
	srv, _ := newTestServer(t)
	resp, body := postSession(t, srv, buildInitData(testAppID, testSecret, nil))
	if resp.StatusCode != 200 || body["auth"] != "guest" || body["user_id"] != nil {
		t.Fatalf("status %d body %v", resp.StatusCode, body)
	}
}

func TestSessionRejectsWithCodeOnly(t *testing.T) {
	srv, _ := newTestServer(t)
	initData := buildInitData(testAppID, "wrong_secret", nil)
	resp, body := postSession(t, srv, initData)
	if resp.StatusCode != 401 || body["error"] != "invalid_signature" || len(body) != 1 {
		t.Fatalf("status %d body %v", resp.StatusCode, body)
	}
}

func TestMeWithoutCookie(t *testing.T) {
	srv, _ := newTestServer(t)
	resp, err := http.Get(srv.URL + "/api/me")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("status %d", resp.StatusCode)
	}
}

func TestMeRestoresSession(t *testing.T) {
	srv, _ := newTestServer(t)
	resp, _ := postSession(t, srv, buildInitData(testAppID, testSecret, Fields{"auth": "customer", "user": testUser}))
	req, _ := http.NewRequest("GET", srv.URL+"/api/me", nil)
	for _, c := range resp.Cookies() {
		req.AddCookie(c)
	}
	me, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	_ = json.NewDecoder(me.Body).Decode(&body)
	if me.StatusCode != 200 || body["user_id"] != "3f1a9c40-77ad-4e0f-9c3d-11c2a0f5e881" {
		t.Fatalf("status %d body %v", me.StatusCode, body)
	}
}

// phoneEnvelope signs an envelope issued "now": the test server checks the
// freshness window against the real clock, same as the session above.
func phoneEnvelope(user string) string {
	issued := time.Now().Unix()
	fields := Fields{
		"app_id":    testAppID,
		"auth_date": strconv.FormatInt(issued, 10),
		"nonce":     "7c9e6679-7425-40de-944b-e07fc1f90ae7",
		"user":      user,
		"phone":     "7011234567",
	}
	fields["sign"] = signFields(fields, "sign", testSecret)
	b, _ := json.Marshal(map[string]any{
		"app_id": fields["app_id"], "auth_date": issued, "nonce": fields["nonce"],
		"user": fields["user"], "phone": fields["phone"], "sign": fields["sign"],
	})
	return string(b)
}

func TestPhoneRequiresSessionAndMatchingUser(t *testing.T) {
	srv, _ := newTestServer(t)

	// no cookie
	resp, _ := http.Post(srv.URL+"/api/phone", "application/json", strings.NewReader(phoneEnvelope(testUser)))
	if resp.StatusCode != 401 {
		t.Fatalf("no session: status %d", resp.StatusCode)
	}

	login, _ := postSession(t, srv, buildInitData(testAppID, testSecret, Fields{"auth": "customer", "user": testUser, "scope": "phone:read"}))

	do := func(env string) (*http.Response, map[string]any) {
		req, _ := http.NewRequest("POST", srv.URL+"/api/phone", strings.NewReader(env))
		req.Header.Set("Content-Type", "application/json")
		for _, c := range login.Cookies() {
			req.AddCookie(c)
		}
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		return r, body
	}

	r, body := do(phoneEnvelope(testUser))
	if r.StatusCode != 200 || body["phone"] != "7011234567" {
		t.Fatalf("valid: status %d body %v", r.StatusCode, body)
	}

	r, body = do(phoneEnvelope(`{"id":"00000000-0000-0000-0000-000000000000"}`))
	if r.StatusCode != 403 || body["error"] != "user_mismatch" {
		t.Fatalf("other user: status %d body %v", r.StatusCode, body)
	}
}
