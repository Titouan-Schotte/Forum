package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/create-account", createAccountHandler)
	http.HandleFunc("/connect-account", connectAccountHandler)
	http.ListenAndServe(":8080", nil)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Check if the pseudo already exists in the database
		email := r.FormValue("email")
		pseudo := r.FormValue("pseudo")
		password := r.FormValue("password")
		fmt.Println(email, pseudo, password)
		http.Redirect(w, r, "/connect-account", http.StatusSeeOther)
	}

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// Il faut ajouter le code pour cree un compte de l'utilisateur
}

func connectAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	// ici le code pour la connection de l'utilisateur
}
