package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

type RegisterData struct {
	ErrorMessage    string
	RegisterSuccess bool
	UserLog         dbmanagement.User
}

var registerData = RegisterData{
	ErrorMessage:    "",
	RegisterSuccess: false,
	UserLog:         dbmanagement.User{},
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	//Request IN
	if r.Method == http.MethodPost {
		pseudoIn := r.FormValue("pseudo")
		emailIn := r.FormValue("email")
		passwordIn := r.FormValue("password")
		confirmPasswordIn := r.FormValue("confirm-password")
		if passwordIn != confirmPasswordIn {
			registerData.ErrorMessage = "Mots de passes différents"
		}
		user, success, errorMsg := dbmanagement.DB.CreateAccount(pseudoIn, emailIn, passwordIn)
		if success {
			registerData.UserLog = user
			registerData.RegisterSuccess = true
		} else {
			registerData.UserLog = dbmanagement.User{}
			registerData.RegisterSuccess = false
			registerData.ErrorMessage = "Error create :" + errorMsg
		}
	}

	tmpl, err := template.ParseFiles("./assets/pages/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, registerData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
