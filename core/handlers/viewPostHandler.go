package handlers

import (
	"Forum/core/dbmanagement"
	"fmt"
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

	// Check if user is logged in, if not redirect to login page
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Get the PostId from the query parameters
	postIDStr := r.URL.Query().Get("PostId")

	if postIDStr == "" {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Convert PostId to integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid PostId", http.StatusBadRequest)
		return
	}

	// Retrieve the post from the database
	post := dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postID)
	viewPostStruct := ViewPost{
		Post:           post,
		IsUserLiked:    loginData.UserLog.IsLikedPost(postID),
		IsUserDisliked: loginData.UserLog.IsDislikedPost(postID),
	}

	for i, comment := range viewPostStruct.Post.Comments {
		fmt.Println(comment.Author.Email)
		fmt.Println(comment.Author.Pseudo)
		author, _, _ := dbmanagement.DB.GetUser(comment.Author.Email)
		viewPostStruct.Post.Comments[i].Author = author
		fmt.Println(viewPostStruct.Post.Comments[i].Author.Pseudo)
	}

	fmt.Println(viewPostStruct.IsUserLiked)

	// Parse the template with the custom function
	tmpl, err := template.ParseFiles("./assets/pages/viewpost.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template using the retrieved post data
	tmpl.Execute(w, viewPostStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
