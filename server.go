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
	"templates/record_delete.html",
	"templates/record_form.html",
	"templates/record_input.html",
	"templates/record_list.html",
))

func renderStringTemplate(w http.ResponseWriter, tmpl string, str string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", str)

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

func recordListViewHandler(w http.ResponseWriter, r *http.Request) {
	records, err := getRecords()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderDataTemplate(w, "record_list", records)
}

func recordDetailsViewHandler(w http.ResponseWriter, r *http.Request) {
	recordId := r.PathValue("recordId")

	record, err := getRecordById(recordId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderDataTemplate(w, "record_details", record)
}

func recordCreateViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle the POST request for creation

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// prepare lastDate from POST string to time.Time
		lastDate, err := time.Parse("2006-01-02", r.FormValue("lastDate"))
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// get POST data
		record := dataRecord{
			Title:    r.FormValue("title"),
			LastDate: lastDate,
			Comment:  r.FormValue("comment"),
		}

		_, createErr := createRecord(record)
		if createErr != nil {
			http.Error(w, "Error creating new record", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	renderDataTemplate(w, "record_input", nil)
}

func recordUpdateViewHandler(w http.ResponseWriter, r *http.Request) {

	recordId := r.PathValue("recordId")

	if r.Method == http.MethodPost {
		// Handle the POST request for update

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// prepare lastDate from POST string to time.Time
		lastDate, err := time.Parse("2006-01-02", r.FormValue("lastDate"))
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// get POST data
		record := dataRecord{
			Id:       recordId,
			Title:    r.FormValue("title"),
			LastDate: lastDate,
			Comment:  r.FormValue("comment"),
		}

		_, createErr := updateRecord(record)
		if createErr != nil {
			http.Error(w, "Error updating record", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	recordToUpdate, err := getRecordById(recordId)
	if err != nil {
		http.Error(w, "Error getting record to update", http.StatusInternalServerError)
		return
	}

	renderDataTemplate(w, "record_input", recordToUpdate)
}

func recordDeleteViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Handle the POST request for deletion

		recordId := r.PathValue("recordId")

		_, err := deleteRecord(recordId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	recordId := r.PathValue("recordId")

	renderStringTemplate(w, "record_delete", recordId)
}

func runServer() {
	port := ":8000"
	router := http.NewServeMux()
	log.Printf("Server starting on %s ...\n", port)

	// static files (CSS)
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	router.HandleFunc("GET /records", recordListViewHandler)
	router.HandleFunc("GET /", recordListViewHandler)

	router.HandleFunc("GET /records/{recordId}", recordDetailsViewHandler)

	router.HandleFunc("GET /records/delete/{recordId}", recordDeleteViewHandler)
	router.HandleFunc("POST /records/delete/{recordId}", recordDeleteViewHandler)

	router.HandleFunc("GET /records/update/", recordCreateViewHandler)
	router.HandleFunc("POST /records/update/", recordCreateViewHandler)

	router.HandleFunc("GET /records/update/{recordId}", recordUpdateViewHandler)
	router.HandleFunc("POST /records/update/{recordId}", recordUpdateViewHandler)

	server := http.Server{
		Addr:    port,
		Handler: loggingMiddleware(router),
	}

	log.Fatal(server.ListenAndServe())
}
