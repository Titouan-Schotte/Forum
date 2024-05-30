package handlers

import (
	"Forum/core/dbmanagement"
	"net/http"
	"strconv"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query parameters
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Call the function to like the post
	postIn := dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postID)
	postIn.LikePost(loginData.UserLog.Email, loginData.UserLog.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the query parameters
	postIDStr := r.URL.Query().Get("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Call the function to dislike the post
	postIn := dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postID)
	postIn.DislikePost(loginData.UserLog.Email, loginData.UserLog.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the comment ID from the query parameters
	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Call the function to like the comment
	commentIn := dbmanagement.DB.GetCommentById(loginData.UserLog.Email, loginData.UserLog.Password, commentID)
	commentIn.LikeComment(loginData.UserLog.Email, loginData.UserLog.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the comment ID from the query parameters
	commentIDStr := r.URL.Query().Get("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Call the function to dislike the comment
	commentIn := dbmanagement.DB.GetCommentById(loginData.UserLog.Email, loginData.UserLog.Password, commentID)
	commentIn.DislikeComment(loginData.UserLog.Email, loginData.UserLog.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}
