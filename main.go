package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/create-account", createAccountHandler)
	http.HandleFunc("/connect-account", connectAccountHandler)
	http.ListenAndServe(":8080", nil)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {
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
