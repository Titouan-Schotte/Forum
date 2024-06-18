package handlers

import (
	"Forum/core/dbmanagement"
	"fmt"
	"html/template"
	"net/http"
)

func PanelGradeHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	//page réservée aux admin !
	if !loginData.UserLog.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//refresh all users
	panelStruct.AllUsers, _ = dbmanagement.DB.GetUsers()
	panelStruct.UserLog = loginData.UserLog
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/grade.html")
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

func PanelActionChangeGradeHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !loginData.UserLog.IsModo && !loginData.UserLog.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/panel", http.StatusSeeOther)
	}
	emailTarget := r.URL.Query().Get("Email")
	gradeIn := r.URL.Query().Get("Grade")
	userTarget, _, _ := dbmanagement.DB.GetUser(emailTarget)

	//Refuse ban !!
	if userTarget.IsAdmin || userTarget.Email == loginData.UserLog.Email {
		tmpl, err := template.ParseFiles("./assets/pages/refuse-changegrade.html")
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
	switch gradeIn {
	case "Normal":
		userTarget.IsModo = false
		userTarget.IsAdmin = false
		break
	case "Administrateur":
		userTarget.IsAdmin = true
	case "Moderateur":
		userTarget.IsModo = true
		break
	}
	fmt.Println(userTarget.IsModo, userTarget.IsAdmin)
	userTarget.EditUser()
	tmpl, err := template.ParseFiles("./assets/pages/success-changegrade.html")
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
