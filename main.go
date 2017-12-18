package main

import (
  "log"
  "net/http"
  "text/template"
  "path/filepath"
  "sync"
  "flag"
  "os"

  "chat/trace"
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
  r := newRoom()
  r.tracer = trace.New(os.Stdout)
  // route
  http.Handle("/", &templateHandler{filename: "chat.html"})
  http.Handle("/room", r)
  // start chat room
  go r.run()
  // start web server
  log.Println("Start Web Server Port: ", *addr)
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
