package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

const cookieName = "sid"

type api struct {
	verifier Verifier
	sessions *sessionStore
	log      *log.Logger
}

func (a *api) routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/session", a.createSession)
	mux.HandleFunc("GET /api/me", a.me)
	mux.HandleFunc("POST /api/phone", a.phone)
}

type sessionResponse struct {
	Auth   string   `json:"auth"`
	UserID *string  `json:"user_id"`
	Scope  []string `json:"scope"`
}

func toResponse(s *session) sessionResponse {
	out := sessionResponse{Auth: s.Auth, Scope: s.Scope}
	if s.PartnerUserID != "" {
		id := s.PartnerUserID
		out.UserID = &id
	}
	if out.Scope == nil {
		out.Scope = []string{}
	}
	return out
}

// createSession takes the launch context as the raw string the page got
// from the app. Parsing it here, not on the client, keeps the bytes the
// signature was computed over.
func (a *api) createSession(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 8<<10))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrMalformed)
		return
	}
	fields, err := parseInitData(string(raw))
	if err != nil {
		a.reject(w, "session", err, fields)
		return
	}
	if err := a.verifier.Verify(fields, "hash"); err != nil {
		a.reject(w, "session", err, fields)
		return
	}

	var scope []string
	if fields["scope"] != "" {
		scope = strings.Split(fields["scope"], ",")
	}

	// #region auth
	var sess *session
	switch fields["auth"] {
	case "customer":
		id, err := userID(fields["user"])
		if err != nil {
			a.reject(w, "session", ErrMalformed, fields)
			return
		}
		// A real backend finds or creates its account keyed by this id.
		sess = a.sessions.create("customer", id, scope)
	case "guest":
		sess = a.sessions.create("guest", "", nil)
	default:
		a.reject(w, "session", ErrMalformed, fields)
		return
	}
	// #endregion auth

	http.SetCookie(w, &http.Cookie{
		Name: cookieName, Value: sess.ID, Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: r.TLS != nil,
	})
	writeJSON(w, http.StatusOK, toResponse(sess))
}

func (a *api) me(w http.ResponseWriter, r *http.Request) {
	sess, ok := a.currentSession(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "no_session"})
		return
	}
	writeJSON(w, http.StatusOK, toResponse(sess))
}

// #region phone

// phone takes the getPhone envelope whole and verifies it like the launch
// context, with sign as the signature field. The user inside must be the
// user of the current session: an envelope is proof about one subscriber,
// and the page could otherwise present someone else's.
func (a *api) phone(w http.ResponseWriter, r *http.Request) {
	sess, ok := a.currentSession(r)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "no_session"})
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 8<<10))
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrMalformed)
		return
	}
	fields, err := parseEnvelope(body)
	if err != nil {
		a.reject(w, "phone", err, fields)
		return
	}
	if err := a.verifier.Verify(fields, "sign"); err != nil {
		a.reject(w, "phone", err, fields)
		return
	}
	id, err := userID(fields["user"])
	if err != nil || id != sess.PartnerUserID {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "user_mismatch"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"phone": fields["phone"]})
}

// #endregion phone

func (a *api) currentSession(r *http.Request) (*session, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil, false
	}
	return a.sessions.get(c.Value)
}

// reject answers 401 with the code only and logs the code with the field
// names present. The context itself never reaches the log: it carries a
// valid signature and the subscriber alias.
func (a *api) reject(w http.ResponseWriter, route string, err error, fields Fields) {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	a.log.Printf("%s rejected: %s fields=%v", route, err, keys)
	writeError(w, http.StatusUnauthorized, err)
}

func userID(user string) (string, error) {
	var parsed struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(user), &parsed); err != nil || parsed.ID == "" {
		return "", errors.New("user without id")
	}
	return parsed.ID, nil
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
