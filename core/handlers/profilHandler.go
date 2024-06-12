package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
)

func ProfilHandler(w http.ResponseWriter, r *http.Request) {

	//Is Not logged in => redirect to login page
	if loginData.UserLog.Email == "" {
		redirectURL := "/login"
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
	if r.URL.Query().Has("email") { //Target profil
		loginData.UserLog, _, _ = dbmanagement.DB.GetUser(r.URL.Query().Get("email"))
		loginData.UserLog.Followers = loginData.UserLog.GetAllFollowers()
		loginData.UserLog.Subscription = loginData.UserLog.GetAllFollowedAccount()
		loginData.UserLog.Posts = loginData.UserLog.GetAllUserPosts()
		loginData.UserLog.Likes = loginData.UserLog.GetAllLikedPosts()
		loginData.UserLog.Notifications = loginData.UserLog.GetAllNotifications()
		loginData.UserLog.TotalLikes = 0
		for _, v := range loginData.UserLog.Posts {
			loginData.UserLog.TotalLikes += v.Like
		}
	} else {
		loginData.UserLog, _, _ = dbmanagement.DB.ConnectToAccount(loginData.UserLog.Email, loginData.UserLog.Password)
	}
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/profil.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute the template using the game data (dataGame)
	err = tmpl.Execute(w, loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
