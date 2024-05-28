package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {

	//Is Not logged in => redirect to login page
	if loginData.UserLog.Email == "" {
		redirectURL := "/login"
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}

	if r.Method == http.MethodGet {
		if r.URL.Query().Get("PostId") == "" {
			PostIn := dbmanagement.DB.GetPostById(r.URL.Query().Get("PostId"))
		}
	}

	loginData.UserLog, _, _ = dbmanagement.DB.ConnectToAccount(loginData.UserLog.Email, loginData.UserLog.Password)
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute the template using the game data (dataGame)
	err = tmpl.Execute(w, loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
