package main

import (
	"github.com/stretchr/gomniauth"
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  "flag"

  "github.com/stretchr/gomniauth/providers/facebook"
  "github.com/stretchr/gomniauth/providers/github"
  "github.com/stretchr/gomniauth/providers/google"  
)

type templateHandler struct {
  once      sync.Once
  filename  string
  templ     *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func(){
    t.templ =
      template.Must(template.ParseFiles(filepath.Join("templates",
        t.filename)))
  })
  t.templ.Execute(w, r)
}

func main() {
  var addr = flag.String("addr", ":8080", "Application Address")
  flag.Parse()

  // setup gomniauth
  gomniauth.SetSecurityKey("hu278guth34gbker72392jhf3yg456sg0148fgbuj2387ugv32hb4o7798ergh2")
  gomniauth.WithProviders(
    facebook.New("208178147561-16rib2nc7hsbgjnde8jtb9k2gm5v5433.apps.googleusercontent.com", "LriuzY2vOdY_HK-UIfPoEHPc", "http://localhost:8080/auth/callback/facebook"),
    github.New("208178147561-65vab8kbmsqbbpa1m5c26g9bpbijjbt4.apps.googleusercontent.com", "RPne_8U_UzUKymZJojSaJIDY", "http://localhost:8080/auth/callback/github"),
    google.New("208178147561-jp8v2fdnf3scet47rhudi6dk2963ebjc.apps.googleusercontent.com", "zc4MlTeQ8qwuCr34OAn2gjIC", "http://localhost:8080/auth/callback/google"),
  )
  r := newRoom()
  // route
  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.HandleFunc("/auth", loginHandler)
  http.Handle("/room", r)
  // start chat room
  go r.run()
  // start web server
  log.Println("Start Web Server Port: ", *addr)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
