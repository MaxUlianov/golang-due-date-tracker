package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

// cache the HTML templates
var templates = template.Must(template.ParseFiles(
	"templates/navbar.html",
	"templates/record_details.html",
	"templates/record_delete.html",
	"templates/record_form.html",
	"templates/record_input.html",
	"templates/record_list.html",
	"templates/user_input.html",
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

// ---- ---- ---- ---- ----
// ---- user auth

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			renderStringTemplate(w, "user_input", "register")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		newUser := userModel{
			Username: username,
			Password: string(hashedPassword),
		}

		_, err = createUser(db_instance, newUser)
		if err != nil {
			http.Error(w, "Error creating user"+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	renderStringTemplate(w, "user_input", "register")
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			renderStringTemplate(w, "user_input", "register")
		}

		userAuthenticated := authUser(db_instance, username, password)
		if !userAuthenticated {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		setSession(w, r, username)

		http.Redirect(w, r, "/records", http.StatusSeeOther)
		return
	}

	renderStringTemplate(w, "user_input", "login")
}

func userLogoutHandler(w http.ResponseWriter, r *http.Request) {
	emptySession(w, r)

	http.Redirect(w, r, "/auth/login/", http.StatusSeeOther)
}

// ---- ---- ---- ---- ----
// ---- session management
func setSession(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := store.Get(r, "session.id")
	session.Values["username"] = username
	session.Save(r, w)
}

func checkSession(r *http.Request) (string, bool) {
	session, _ := store.Get(r, "session.id")
	username, ok := session.Values["username"].(string)
	return username, ok
}

func emptySession(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	delete(session.Values, "username")
	session.Save(r, w)
}

// ---- ---- ---- ---- ----
// ---- record Handlers

func recordListViewHandler(w http.ResponseWriter, r *http.Request) {
	username, _ := checkSession(r)
	userId, err := getUserId(db_instance, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	records, err := getRecords(db_instance, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderDataTemplate(w, "record_list", records)
}

func recordDetailsViewHandler(w http.ResponseWriter, r *http.Request) {
	recordId := r.PathValue("recordId")

	username, _ := checkSession(r)
	userId, err := getUserId(db_instance, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	record, err := getRecordById(db_instance, recordId, userId)
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

		username, _ := checkSession(r)
		userId, err := getUserId(db_instance, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, createErr := createRecord(db_instance, record, userId)
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

		username, _ := checkSession(r)
		userId, err := getUserId(db_instance, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, createErr := updateRecord(db_instance, record, userId)
		if createErr != nil {
			http.Error(w, "Error updating record", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username, _ := checkSession(r)
	userId, err := getUserId(db_instance, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recordToUpdate, err := getRecordById(db_instance, recordId, userId)
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

		username, _ := checkSession(r)
		userId, err := getUserId(db_instance, username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, deleteErr := deleteRecord(db_instance, recordId, userId)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
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

	router.HandleFunc("GET /auth/register/", userRegisterHandler)
	router.HandleFunc("POST /auth/register/", userRegisterHandler)

	router.HandleFunc("GET /auth/login/", userLoginHandler)
	router.HandleFunc("POST /auth/login/", userLoginHandler)

	router.HandleFunc("GET /auth/logout/", userLogoutHandler)

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
		Handler: loggingMiddleware(authMiddleware(router)),
	}

	log.Fatal(server.ListenAndServe())
}
