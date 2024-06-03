package handlers

import (
	"Forum/core/dbmanagement"
	"net/http"
)

func DisconnectHandler(w http.ResponseWriter, r *http.Request) {
	//Is Not logged in => redirect to login page

	loginData = LoginData{
		ErrorMessage:    "",
		RegisterSuccess: false,
		UserLog:         dbmanagement.User{},
	}
	// Load the home page template
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
