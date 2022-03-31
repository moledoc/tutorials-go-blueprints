package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenicated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// success - call next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth is creates authHandler that wraps any other http.Handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginhandler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 4 || segs[1] != "auth" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Provided path (%s) not supported", r.URL.Path)
		return
	}
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Println("TODO: handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
