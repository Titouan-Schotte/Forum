/*
Titouan Schotté

Panel grade changer handler
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func PanelGradeHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !loginData.UserLog.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	panelStruct.AllUsers, _ = dbmanagement.DB.GetUsers()
	panelStruct.UserLog = loginData.UserLog
	tmpl, err := template.ParseFiles("./assets/pages/grade.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	if userTarget.IsAdmin || userTarget.Email == loginData.UserLog.Email {
		tmpl, err := template.ParseFiles("./assets/pages/refuse-changegrade.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
	userTarget.EditUser()
	tmpl, err := template.ParseFiles("./assets/pages/success-changegrade.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, panelStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
