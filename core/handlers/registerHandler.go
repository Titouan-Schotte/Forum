/*
Titouan Schotté

Register new user handler
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.FormValue("pseudo") != "" {
			pseudoIn := r.FormValue("pseudo")
			emailIn := r.FormValue("email")
			passwordIn := r.FormValue("password")
			confirmPasswordIn := r.FormValue("confirm-password")

			if passwordIn != confirmPasswordIn {
				loginData.ErrorMessage = "Mots de passes différents"
			} else {
				user, success, errorMsg := dbmanagement.DB.CreateAccount(pseudoIn, emailIn, passwordIn)
				if success {
					loginData.UserLog = user
					loginData.RegisterSuccess = true

					redirectURL := "/login?email=" + user.Email + "&password=" + user.Password + "&registerSuccess=true"
					http.Redirect(w, r, redirectURL, http.StatusSeeOther)
					return
				} else {
					loginData.UserLog = dbmanagement.User{}
					loginData.RegisterSuccess = false
					loginData.ErrorMessage = "Error create :" + errorMsg
				}
			}
		}
	}

	tmpl, err := template.ParseFiles("./assets/pages/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
