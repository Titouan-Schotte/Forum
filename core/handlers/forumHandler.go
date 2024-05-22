package handlers

import "net/http"
import "html/template"

type CoreDatas struct {
}

var coreDatas = CoreDatas{}

func ForumHandler(w http.ResponseWriter, r *http.Request) {
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/forum.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute the template using the game data (dataGame)
	err = tmpl.Execute(w, coreDatas)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
