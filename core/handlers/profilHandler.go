package handlers

import (
	"html/template"
	"net/http"
)

func ProfilHandler(w http.ResponseWriter, r *http.Request) {
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/profil.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute the template using the game data (dataGame)
	err = tmpl.Execute(w, coreDatas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
