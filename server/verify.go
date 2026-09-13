package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Rejection codes. The text is what the client sees in the 401 body; the
// reason behind it stays in the server log.
var (
	ErrInvalidSignature = errors.New("invalid_signature")
	ErrExpired          = errors.New("expired")
	ErrReplayed         = errors.New("replayed")
	ErrAppMismatch      = errors.New("app_mismatch")
	ErrMalformed        = errors.New("malformed")
)

// Fields is a parsed launch context or getPhone envelope: field name to
// value exactly as decoded, before any interpretation.
type Fields map[string]string

// parseInitData parses the launch context query string. url.ParseQuery
// percent-decodes the values, so user comes back as plain JSON text.
func parseInitData(raw string) (Fields, error) {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil, ErrMalformed
	}
	fields := make(Fields, len(values))
	for key, list := range values {
		if len(list) != 1 {
			return nil, ErrMalformed
		}
		fields[key] = list[0]
	}
	return fields, nil
}

// parseEnvelope parses the getPhone envelope. auth_date arrives as a JSON
// number and goes into the canonical string in its decimal form; user is a
// JSON string and must stay byte-for-byte as received.
func parseEnvelope(body []byte) (Fields, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, ErrMalformed
	}
	fields := make(Fields, len(raw))
	for key, value := range raw {
		var s string
		if err := json.Unmarshal(value, &s); err == nil {
			fields[key] = s
			continue
		}
		var n json.Number
		if err := json.Unmarshal(value, &n); err == nil {
			fields[key] = n.String()
			continue
		}
		return nil, ErrMalformed
	}
	return fields, nil
}

// canonical builds the signed string: key=value pairs without the signature
// field, keys sorted, joined with "\n", no trailing newline.
func canonical(fields Fields, signKey string) string {
	keys := make([]string, 0, len(fields))
	for key := range fields {
		if key != signKey {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	pairs := make([]string, len(keys))
	for i, key := range keys {
		pairs[i] = key + "=" + fields[key]
	}
	return strings.Join(pairs, "\n")
}

func checkSignature(fields Fields, signKey, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonical(fields, signKey)))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(fields[signKey]))
}

// nonceStore remembers presented nonces for a TTL. It lives in process
// memory, which is only correct for a single instance: with two replicas a
// replay can land on the one that has not seen the nonce. That is the cost
// that made the check optional; a real deployment needs Redis or a table.
type nonceStore struct {
	mu   sync.Mutex
	ttl  time.Duration
	seen map[string]time.Time
}

func newNonceStore(ttl time.Duration) *nonceStore {
	return &nonceStore{ttl: ttl, seen: map[string]time.Time{}}
}

// Add records the nonce and reports whether it was new.
func (s *nonceStore) Add(nonce string, now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, at := range s.seen {
		if now.Sub(at) > s.ttl {
			delete(s.seen, key)
		}
	}
	if at, ok := s.seen[nonce]; ok && now.Sub(at) <= s.ttl {
		return false
	}
	s.seen[nonce] = now
	return true
}

// Verifier checks a context or envelope for one mini app.
type Verifier struct {
	AppID  string
	Secret string
	// Window is the accepted skew of auth_date, both directions.
	Window time.Duration
	// Nonces is nil when NONCE_CHECK=off.
	Nonces *nonceStore
	Now    func() time.Time
}

// Verify runs the checks in the order the guide lists them. The signature
// goes first: until it holds, no field means anything.
func (v Verifier) Verify(fields Fields, signKey string) error {
	if !checkSignature(fields, signKey, v.Secret) {
		return ErrInvalidSignature
	}
	authDate, err := strconv.ParseInt(fields["auth_date"], 10, 64)
	if err != nil {
		return ErrMalformed
	}
	skew := v.Now().Unix() - authDate
	if skew < 0 {
		skew = -skew
	}
	if skew > int64(v.Window/time.Second) {
		return ErrExpired
	}
	if v.Nonces != nil && !v.Nonces.Add(fields["nonce"], v.Now()) {
		return ErrReplayed
	}
	if fields["app_id"] != v.AppID {
		return ErrAppMismatch
	}
	return nil
}
