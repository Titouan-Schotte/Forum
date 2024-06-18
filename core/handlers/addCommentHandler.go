package handlers

import (
	"Forum/core/dbmanagement"
	"html/template"
	"net/http"
	"strconv"
)

func PanelAddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !loginData.UserLog.IsModo && !loginData.UserLog.IsAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	postId, _ := strconv.Atoi(r.URL.Query().Get("PostId"))
	content := r.URL.Query().Get("Content")
	//Refuse ban !!
	if loginData.UserLog.IsBan {
		tmpl, err := template.ParseFiles("./assets/pages/isban.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Execute the template using the game data (dataGame)
		err = tmpl.Execute(w, panelStruct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	postIn := dbmanagement.DB.GetPostById(loginData.UserLog.Email, loginData.UserLog.Password, postId)
	postIn.Author.AddNotification("Votre post "+r.FormValue("title")+" a été commenté par "+loginData.UserLog.Pseudo, "comment")

	postIn.AddComment(loginData.UserLog.Email, loginData.UserLog.Password, content)
	http.Redirect(w, r, "/viewpost?PostId="+r.URL.Query().Get("PostId"), http.StatusSeeOther)

}
