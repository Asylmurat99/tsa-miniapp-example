package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

// #region check
// checkSignature verifies a launch context (signField "hash") or a getPhone
// envelope (signField "sign"). fields holds decoded values.
func checkSignature(fields map[string]string, secret, signField string) bool {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		if k != signField {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	pairs := make([]string, len(keys))
	for i, k := range keys {
		pairs[i] = k + "=" + fields[k]
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(strings.Join(pairs, "\n")))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(fields[signField]))
}
// #endregion check

func main() {
	data, err := os.ReadFile("../vector.json")
	if err != nil {
		panic(err)
	}
	var v struct {
		Secret string `json:"secret"`
		Cases  []struct {
			Name      string            `json:"name"`
			SignField string            `json:"sign_field"`
			Fields    map[string]string `json:"fields"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		panic(err)
	}
	failed := false
	for _, c := range v.Cases {
		if checkSignature(c.Fields, v.Secret, c.SignField) {
			fmt.Printf("%s: ok\n", c.Name)
		} else {
			fmt.Printf("%s: FAIL\n", c.Name)
			failed = true
		}
	}
	if failed {
		os.Exit(1)
	}
}
