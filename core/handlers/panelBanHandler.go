/*
Titouan Schotté

Panel ban handler
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func PanelBanHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !loginData.UserLog.IsModo && !loginData.UserLog.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	panelStruct.AllUsers, _ = dbmanagement.DB.GetUsers()
	panelStruct.UserLog = loginData.UserLog
	tmpl, err := template.ParseFiles("./assets/pages/bannis.html")
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

func PanelActionBanHandler(w http.ResponseWriter, r *http.Request) {
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
	userTarget, _, _ := dbmanagement.DB.GetUser(emailTarget)
	if userTarget.IsModo || userTarget.IsAdmin || userTarget.Email == loginData.UserLog.Email {
		tmpl, err := template.ParseFiles("./assets/pages/refuse-ban.html")
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

	panelStruct.UserLog.BanAccountModo(emailTarget)
	tmpl, err := template.ParseFiles("./assets/pages/success-ban.html")
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
