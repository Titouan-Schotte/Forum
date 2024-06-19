/*
Titouan Schotté

Settings page handler & actions handlers
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/pages/settings.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, forumPostsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SettingsPasswordChange(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		currentPassword := r.FormValue("current-password")
		newPassword := r.FormValue("new-password")
		confirmPassword := r.FormValue("confirm-new-password")

		if newPassword != confirmPassword {
			http.Redirect(w, r, "/settings?errormessage=Mots de passes actuel différents.", http.StatusSeeOther)
			return
		}
		user, authenticated := dbmanagement.IsUserConnected(loginData.UserLog.Email, currentPassword)
		if !authenticated {
			http.Redirect(w, r, "/settings?errormessage=Mot de passe actuel incorrect.", http.StatusSeeOther)

			return
		}

		user.Password = newPassword
		user.EditUser()
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
	}
}

func SettingsPseudoChange(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		newPseudo := r.FormValue("new-username")
		user, authenticated := dbmanagement.IsUserConnected(loginData.UserLog.Email, loginData.UserLog.Password)
		if !authenticated {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user.Pseudo = newPseudo
		user.EditUser()

		http.Redirect(w, r, "/settings", http.StatusSeeOther)
	}
}

func SettingsDeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		user, authenticated := dbmanagement.IsUserConnected(loginData.UserLog.Email, loginData.UserLog.Password)
		if !authenticated {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		success, errMsg := user.DeleteAccount()
		if !success {
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
