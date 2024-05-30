package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

type PanelStruct struct {
	AllUsers []dbmanagement.User
	UserLog  dbmanagement.User
}

var panelStruct = PanelStruct{}

func PanelAdminHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !loginData.UserLog.IsModo && !loginData.UserLog.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//Dump all users

	panelStruct.AllUsers, _ = dbmanagement.DB.GetUsers()
	panelStruct.UserLog = loginData.UserLog
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/admin.html")
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
}
