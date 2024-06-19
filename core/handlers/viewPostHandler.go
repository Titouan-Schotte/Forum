/*
Titouan Schotté

View post details handlers
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
	"strconv"
)

type ViewPost struct {
	Post           dbmanagement.Post
	IsUserLiked    bool
	IsUserDisliked bool
}

func ViewPostHandler(w http.ResponseWriter, r *http.Request) {
	redirectURL := "/login"

	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	postIDStr := r.URL.Query().Get("PostId")

	if postIDStr == "" {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid PostId", http.StatusBadRequest)
		return
	}

	post := dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postID)
	viewPostStruct := ViewPost{
		Post:           post,
		IsUserLiked:    loginData.UserLog.IsLikedPost(postID),
		IsUserDisliked: loginData.UserLog.IsDislikedPost(postID),
	}

	for i, comment := range viewPostStruct.Post.Comments {
		author, _, _ := dbmanagement.DB.GetUser(comment.Author.Email)
		viewPostStruct.Post.Comments[i].Author = author
	}

	tmpl, err := template.ParseFiles("./assets/pages/viewpost.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, viewPostStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
