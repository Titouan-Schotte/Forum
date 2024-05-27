package handlers

import (
	"Forum/core/dbmanagement"
	"net/http"
)
import "html/template"

type ForumPostsGetter struct {
	RecentPosts []dbmanagement.Post
	TopPosts    []dbmanagement.Post
	UnePosts    []dbmanagement.Post
}

var forumPostsData = ForumPostsGetter{}

func ForumHandler(w http.ResponseWriter, r *http.Request) {

	//LOGIN IN !!
	if loginData.UserLog.Email != "" {

	}

	forumPostsData.RecentPosts = dbmanagement.DB.GetMostRecentsPosts(10)
	forumPostsData.TopPosts = dbmanagement.DB.GetTopPosts(10)
	forumPostsData.UnePosts = dbmanagement.DB.GetRandomPosts(10)

	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/forum.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute the template using the game data (dataGame)
	err = tmpl.Execute(w, forumPostsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
