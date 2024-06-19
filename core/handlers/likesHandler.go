/*
Titouan Schotté

Likes & Dislikes actions handlers
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"net/http"
	"strconv"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	postIn := dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postID)
	postIn.Author.AddNotification("Votre post "+r.FormValue("title")+" a été liké par "+loginData.UserLog.Pseudo, "like")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postIn := dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postID)
	postIn.Author.AddNotification("Votre post "+r.FormValue("title")+" a été disliké par "+loginData.UserLog.Pseudo, "dislike")

	postIn.DislikePost(loginData.UserLog.Email, loginData.UserLog.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
