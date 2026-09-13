package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

type vectorFile struct {
	Secret string `json:"secret"`
	Cases  []struct {
		Name      string            `json:"name"`
		SignField string            `json:"sign_field"`
		Fields    map[string]string `json:"fields"`
		Raw       string            `json:"raw"`
		JSON      string            `json:"json"`
	} `json:"cases"`
}

func loadVectors(t *testing.T) vectorFile {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "verify", "vector.json"))
	if err != nil {
		t.Fatal(err)
	}
	var v vectorFile
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestSignatureMatchesVectors(t *testing.T) {
	v := loadVectors(t)
	for _, c := range v.Cases {
		if !checkSignature(c.Fields, c.SignField, v.Secret) {
			t.Errorf("%s: signature rejected", c.Name)
		}
		tampered := Fields{}
		for k, val := range c.Fields {
			tampered[k] = val
		}
		tampered["auth_date"] = "1"
		if checkSignature(tampered, c.SignField, v.Secret) {
			t.Errorf("%s: tampered fields accepted", c.Name)
		}
	}
}

func TestParseInitDataDecodesValues(t *testing.T) {
	v := loadVectors(t)
	for _, c := range v.Cases {
		if c.Raw == "" {
			continue
		}
		fields, err := parseInitData(c.Raw)
		if err != nil {
			t.Fatalf("%s: %v", c.Name, err)
		}
		for k, want := range c.Fields {
			if fields[k] != want {
				t.Errorf("%s: %s = %q, want %q", c.Name, k, fields[k], want)
			}
		}
		if len(fields) != len(c.Fields) {
			t.Errorf("%s: %d fields, want %d", c.Name, len(fields), len(c.Fields))
		}
	}
}

func TestParseEnvelopeKeepsNumbersAsDecimalStrings(t *testing.T) {
	v := loadVectors(t)
	for _, c := range v.Cases {
		if c.JSON == "" {
			continue
		}
		fields, err := parseEnvelope([]byte(c.JSON))
		if err != nil {
			t.Fatalf("%s: %v", c.Name, err)
		}
		for k, want := range c.Fields {
			if fields[k] != want {
				t.Errorf("%s: %s = %q, want %q", c.Name, k, fields[k], want)
			}
		}
	}
}

func TestCanonicalSortsKeysAndSkipsSignature(t *testing.T) {
	got := canonical(Fields{"b": "2", "hash": "x", "a": "1"}, "hash")
	if got != "a=1\nb=2" {
		t.Fatalf("canonical = %q", got)
	}
}

func verifierFor(t *testing.T, v vectorFile, authDate string, nonces *nonceStore) Verifier {
	t.Helper()
	issued, err := strconv.ParseInt(authDate, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return Verifier{
		AppID:  "34d5810c-7c36-11eb-82e7-f2189812cd57",
		Secret: v.Secret,
		Window: 300 * time.Second,
		Nonces: nonces,
		Now:    func() time.Time { return time.Unix(issued+10, 0) },
	}
}

func TestVerifyOrder(t *testing.T) {
	v := loadVectors(t)
	c := v.Cases[0]

	t.Run("valid", func(t *testing.T) {
		if err := verifierFor(t, v, c.Fields["auth_date"], nil).Verify(c.Fields, c.SignField); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("signature checked before anything else", func(t *testing.T) {
		vf := verifierFor(t, v, c.Fields["auth_date"], nil)
		vf.AppID = "other"
		bad := copyFields(c.Fields)
		bad["hash"] = "00"
		if err := vf.Verify(bad, c.SignField); !errors.Is(err, ErrInvalidSignature) {
			t.Fatalf("err = %v, want invalid_signature", err)
		}
	})

	t.Run("expired both directions", func(t *testing.T) {
		vf := verifierFor(t, v, c.Fields["auth_date"], nil)
		vf.Now = func() time.Time { return time.Unix(1755388800+301, 0) }
		if err := vf.Verify(c.Fields, c.SignField); !errors.Is(err, ErrExpired) {
			t.Fatalf("future now: err = %v", err)
		}
		vf.Now = func() time.Time { return time.Unix(1755388800-301, 0) }
		if err := vf.Verify(c.Fields, c.SignField); !errors.Is(err, ErrExpired) {
			t.Fatalf("past now: err = %v", err)
		}
	})

	t.Run("replay only when store is on", func(t *testing.T) {
		vf := verifierFor(t, v, c.Fields["auth_date"], newNonceStore(900*time.Second))
		if err := vf.Verify(c.Fields, c.SignField); err != nil {
			t.Fatal(err)
		}
		if err := vf.Verify(c.Fields, c.SignField); !errors.Is(err, ErrReplayed) {
			t.Fatalf("second use: err = %v", err)
		}
		off := verifierFor(t, v, c.Fields["auth_date"], nil)
		if err := off.Verify(c.Fields, c.SignField); err != nil {
			t.Fatalf("store off must accept repeat: %v", err)
		}
	})

	t.Run("app mismatch", func(t *testing.T) {
		vf := verifierFor(t, v, c.Fields["auth_date"], nil)
		vf.AppID = "00000000-0000-0000-0000-000000000000"
		if err := vf.Verify(c.Fields, c.SignField); !errors.Is(err, ErrAppMismatch) {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("malformed auth_date", func(t *testing.T) {
		bad := copyFields(c.Fields)
		bad["auth_date"] = "yesterday"
		bad["hash"] = signFields(bad, "hash", v.Secret)
		if err := verifierFor(t, v, c.Fields["auth_date"], nil).Verify(bad, c.SignField); !errors.Is(err, ErrMalformed) {
			t.Fatalf("err = %v", err)
		}
	})
}

func TestNonceStoreForgetsAfterTTL(t *testing.T) {
	s := newNonceStore(900 * time.Second)
	t0 := time.Unix(1000, 0)
	if !s.Add("n", t0) {
		t.Fatal("first add rejected")
	}
	if s.Add("n", t0.Add(899*time.Second)) {
		t.Fatal("repeat inside ttl accepted")
	}
	if !s.Add("n", t0.Add(901*time.Second)) {
		t.Fatal("repeat after ttl rejected")
	}
}

func copyFields(src Fields) Fields {
	dst := make(Fields, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
