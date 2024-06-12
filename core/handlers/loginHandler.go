package handlers

import (
	"Forum/core/dbmanagement"
	"fmt"
	"html/template"
	"net/http"
)

type LoginData struct {
	ErrorMessage        string
	RegisterSuccess     bool
	UserLog             dbmanagement.User
	CategoriesAvailable []dbmanagement.Categorie
}

var loginData = LoginData{
	ErrorMessage:    "",
	RegisterSuccess: false,
	UserLog:         dbmanagement.User{},
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		fmt.Print("QUOICOU")
		var emailIn string
		var passwordIn string
		if r.Method == http.MethodPost {
			emailIn = r.FormValue("email")
			passwordIn = r.FormValue("password")
		}
		if r.URL.Query().Get("email") != "" {
			emailIn = r.URL.Query().Get("email")
			passwordIn = r.URL.Query().Get("password")
		}

		user, success, errorMsg := dbmanagement.DB.ConnectToAccount(emailIn, passwordIn)
		if success {
			loginData.UserLog = user
			loginData.RegisterSuccess = true
		} else if emailIn != "" {
			fmt.Println(emailIn)
			loginData.UserLog = dbmanagement.User{}
			loginData.RegisterSuccess = false
			loginData.ErrorMessage = "Error connect: " + errorMsg
		} else {
			loginData = LoginData{
				ErrorMessage:    "",
				RegisterSuccess: false,
				UserLog:         dbmanagement.User{},
			}
		}
	}

	tmpl, err := template.ParseFiles("./assets/pages/login.html")
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
