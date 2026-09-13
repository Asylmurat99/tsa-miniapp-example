package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// signFields returns the HMAC-SHA256 of the canonical string, lowercase hex.
// The server never signs anything in normal operation; this exists for the
// -sign flag and for tests, so that a context can be produced without the
// mobile app.
func signFields(fields Fields, signKey, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(canonical(fields, signKey)))
	return hex.EncodeToString(mac.Sum(nil))
}

// buildInitData assembles a launch context query string the way the platform
// does: current auth_date, fresh nonce, given app_id, hash last.
func buildInitData(appID, secret string, extra Fields) string {
	fields := Fields{
		"app_id":    appID,
		"auth_date": strconv.FormatInt(time.Now().Unix(), 10),
		"nonce":     randomUUID(),
		"platform":  "android",
		"auth":      "guest",
	}
	for key, value := range extra {
		fields[key] = value
	}
	values := url.Values{}
	for key, value := range fields {
		values.Set(key, value)
	}
	return values.Encode() + "&hash=" + signFields(fields, "hash", secret)
}

func randomUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
