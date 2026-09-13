package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	signGuest := flag.Bool("sign", false, "print a freshly signed launch context (guest by default) and exit")
	signUser := flag.String("sign-user", "", "with -sign: partner user id to embed as a customer context")
	signScope := flag.String("sign-scope", "", "with -sign-user: scope list, e.g. phone:read")
	flag.Parse()

	appID := mustEnv("APP_ID")
	secret := mustEnv("APP_SECRET")

	if *signGuest {
		extra := Fields{}
		if *signUser != "" {
			extra["auth"] = "customer"
			extra["user"] = fmt.Sprintf(`{"id":%q}`, *signUser)
			if *signScope != "" {
				extra["scope"] = *signScope
			}
		}
		fmt.Println(buildInitData(appID, secret, extra))
		return
	}

	window, err := strconv.Atoi(envOr("AUTH_WINDOW", "300"))
	if err != nil {
		log.Fatalf("AUTH_WINDOW: %v", err)
	}

	verifier := Verifier{AppID: appID, Secret: secret, Window: time.Duration(window) * time.Second, Now: time.Now}
	if envOr("NONCE_CHECK", "off") == "on" {
		// 900 s = the 600 s acceptance interval plus margin, see the guide.
		verifier.Nonces = newNonceStore(900 * time.Second)
	}

	a := &api{verifier: verifier, sessions: newSessionStore(), log: log.Default()}
	mux := http.NewServeMux()
	a.routes(mux)
	// The documentation site is the front door; the mini app itself lives
	// under /app/ and is what the Telecom app opens.
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("web/dist"))))
	mux.Handle("/", http.FileServer(http.Dir("docs/.vitepress/dist")))

	addr := envOr("ADDR", ":8080")
	log.Printf("listening on %s, nonce check %s", addr, envOr("NONCE_CHECK", "off"))
	log.Fatal(http.ListenAndServe(addr, mux))
}

func mustEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("%s is not set; copy .env.example to .env and load it", name)
	}
	return v
}

func envOr(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
