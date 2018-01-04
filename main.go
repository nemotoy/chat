package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
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

var host = flag.String("host", ":8080", "The host of the application.")

func main() {

	flag.Parse() // parse the flags

	// setup gomniauth
	gomniauth.SetSecurityKey("98dfbg7iu2nb4uywevihjw4tuiyub34noilk")
	gomniauth.WithProviders(
		github.New("65vab8kbmsqbbpa1m5c26g9bpbijjbt4", "RPne_8U_UzUKymZJojSaJIDY", "http://localhost:8080/auth/callback/github"),
		google.New("208178147561-jp8v2fdnf3scet47rhudi6dk2963ebjc.apps.googleusercontent.com", "zc4MlTeQ8qwuCr34OAn2gjIC", "http://localhost:8080/auth/callback/google"),
		facebook.New("16rib2nc7hsbgjnde8jtb9k2gm5v5433", "LriuzY2vOdY_HK-UIfPoEHPc", "http://localhost:8080/auth/callback/facebook"),
	)

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}