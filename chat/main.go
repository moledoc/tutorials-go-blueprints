package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/moledoc/tutorials-go-blueprints/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// set the active Avatar implementation
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	addr := flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	// setup gomniauth
	gomniauth.SetSecurityKey("This is a security phrase (which is not so good maybe)")
	gomniauth.WithProviders(
		google.New(os.Getenv("AUTH_ID"), os.Getenv("AUTH_KEY"), "http://localhost:8080/auth/callback/google"),
		// TODO: facebook.New(os.Getenv("AUTH_ID"), os.Getenv("AUTH_KEY"), "http://localhost:8080/auth/callback/facebook"),
		// TODO: github.New(os.Getenv("AUTH_ID"), os.Getenv("AUTH_KEY"), "http://localhost:8080/auth/callback/github"),
	)
	err := os.MkdirAll("./avatars", os.ModePerm)
	if err != nil {
		log.Println("Could not create ./avatars directory")
		return
	}
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	go r.run()
	log.Printf("Serving %v\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
