//go 1.10.4

package main

import(
  "os"
  "strconv"
  "log"
  "net/http"
  "net/smtp"
  "github.com/gorilla/mux"
  "html/template"
  "encoding/json"
  )

var templates *template.Template

type contactForm struct {
  Email string "email"
  Comment string "comment"
}

func main(){

  port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    tStr := os.Getenv("REPEAT")
    repeat, err := strconv.Atoi(tStr)
    if err != nil {
        log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
        repeat = 5
    }
  
  templates = template.Must(template.ParseGlob("templates/*.html"))
  
  r := mux.NewRouter()
  r.HandleFunc("/", ResumeHandler).Methods("GET")
  r.HandleFunc("/Send", Send).Methods("POST")
  r.HandleFunc("/Confirmation", Confirmation).Methods("GET")
  
  fs := http.FileServer(http.Dir("./static/"))
  r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))
  
  http.Handle("/", r)
  
  log.Println("Server is Listening...")
  log.Fatal(http.ListenAndServe(":" + port, r))
}


func ResumeHandler(w http.ResponseWriter, r *http.Request) {
  templates.ExecuteTemplate(w, "index.html", r)
}


func Send(w http.ResponseWriter, r *http.Request) {
  c := &contactForm{}
  json.NewDecoder(r.Body).Decode(c)
  
  to := "gerepromise@gmail.com"
  subject := "NEW CONTACT"
  body := "To: " + to + "\r\nSubject: " + subject + "\r\n\r\n" + "Email: " + c.Email + "\r\n\r\n" + "Comment: " + c.Comment
  auth := smtp.PlainAuth("", "gerepromise@gmail.com", "Imadeit55", "smtp.gmail.com")
  err := smtp.SendMail("smtp.gmail.com:587", auth, "gerepromise@gmail.com", []string{to},[]byte(body))
  if err != nil {
    log.Println("attempting to send mail", err)
  }
  http.Redirect(w, r, "/Confirmation", 302)
}


func Confirmation(w http.ResponseWriter, r *http.Request) {
  templates.ExecuteTemplate(w, "confirm.html", r)
}


