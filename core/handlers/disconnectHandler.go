/*
Titouan Schotté

Disconnect handler
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"net/http"
)

func DisconnectHandler(w http.ResponseWriter, r *http.Request) {

	loginData = LoginData{
		ErrorMessage:    "",
		RegisterSuccess: false,
		UserLog:         dbmanagement.User{},
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
