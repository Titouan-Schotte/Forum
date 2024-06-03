package handlers

import (
	"Forum/core/dbmanagement"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type ViewPost struct {
	Post                       dbmanagement.Post
	IsUserLiked                bool
	IsUserDisliked             bool
	IsUserLikedCommentaries    []bool
	IsUserDislikedCommentaries []bool
}

func toJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
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
	viewPostStruct := ViewPost{
		Post:                       dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postID),
		IsUserLiked:                loginData.UserLog.IsLikedPost(postID),
		IsUserDisliked:             loginData.UserLog.IsDislikedPost(postID),
		IsUserLikedCommentaries:    []bool{},
		IsUserDislikedCommentaries: []bool{},
	}
	for i, comment := range viewPostStruct.Post.Comments {
		viewPostStruct.IsUserLikedCommentaries = append(viewPostStruct.IsUserLikedCommentaries, loginData.UserLog.IsLikedComment(comment.Id))
		viewPostStruct.IsUserDislikedCommentaries = append(viewPostStruct.IsUserDislikedCommentaries, loginData.UserLog.IsDislikedComment(comment.Id))
		fmt.Println(comment.Author.Email)
		fmt.Println(comment.Author.Pseudo)
		viewPostStruct.Post.Comments[i].Author, _, _ = dbmanagement.DB.GetUser(comment.Author.Email)
		fmt.Println(viewPostStruct.Post.Comments[i].Author.Pseudo)

	}
	fmt.Println(viewPostStruct.IsUserLiked)
	// Parse the template with the custom function
	tmpl, err := template.New("viewpost.html").Funcs(template.FuncMap{
		"toJSON": toJSON,
	}).ParseFiles("./assets/pages/viewpost.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template using the retrieved post data
	err = tmpl.Execute(w, viewPostStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
