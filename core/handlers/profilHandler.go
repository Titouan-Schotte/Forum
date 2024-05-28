package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	//Is Not logged in => redirect to login page
	if loginData.UserLog.Email == "" {
		redirectURL := "/login"
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
	loginData.UserLog, _, _ = dbmanagement.DB.ConnectToAccount(loginData.UserLog.Email, loginData.UserLog.Password)
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/profil.html")
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
