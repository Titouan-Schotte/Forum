package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func PanelAddCatHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !loginData.UserLog.IsModo && !loginData.UserLog.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		categorieName := r.FormValue("categorieName")
		dbmanagement.DB.AddCategorie(loginData.UserLog.Email, loginData.UserLog.Password, categorieName)
		tmpl, err := template.ParseFiles("./assets/pages/success-addcat.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Execute the template using the game data (dataGame)
		err = tmpl.Execute(w, panelStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/addcategories.html")
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
