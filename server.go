package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

// cache the HTML templates
var templates = template.Must(template.ParseFiles(
	"templates/navbar.html",
	"templates/record_details.html",
	"templates/record_form.html",
	"templates/record_list.html",
))

func renderTemplate(w http.ResponseWriter, tmpl string, title string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", title)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderDataTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "record_list", "")
}

var mockRecords = []dataRecord{
	{
		Id:       "0000",
		Title:    "Record 1",
		Comment:  "This is the first record.",
		LastDate: time.Now().AddDate(0, 0, -7),
	},
	{
		Id:       "0000",
		Title:    "Record 2",
		Comment:  "This is the second record.",
		LastDate: time.Now().AddDate(0, 0, -3),
	},
	{
		Id:       "0000",
		Title:    "Record 3",
		Comment:  "This is the third record.",
		LastDate: time.Now(),
	},
}

func listViewHandler(w http.ResponseWriter, r *http.Request) {
	renderDataTemplate(w, "record_list", mockRecords)
}

// var validPath = regexp.MustCompile("^/web/([a-zA-Z0-9]+)$")

// func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
//         return func(w http.ResponseWriter, r *http.Request) {
//                 m := validPath.FindStringSubmatch(r.URL.Path)
//                 if m == nil {
//                         http.NotFound(w, r)
//                         return
//                 }

//                 fn(w, r, m[1])
//         }
// }

func runServer() {
	port := ":8000"
	router := http.NewServeMux()
	log.Printf("Server starting on %s ...\n", port)

	// static files (CSS)
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.HandleFunc("GET /records", listViewHandler)

	server := http.Server{
		Addr:    port,
		Handler: loggingMiddleware(router),
	}

	log.Fatal(server.ListenAndServe())
}
