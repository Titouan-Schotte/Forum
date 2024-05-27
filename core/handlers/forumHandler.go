package handlers

import (
	"net/http"
)
import "html/template"

func ForumHandler(w http.ResponseWriter, r *http.Request) {

	//LOGIN IN !!
	if loginData.UserLog.Email != "" {

	}

	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/forum.html")
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
