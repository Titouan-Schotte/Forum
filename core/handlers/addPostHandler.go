/*
Titouan Schotté

Add post handler
*/
package handlers

import (
	"Forum/core/dbmanagement"
	"github.com/gofrs/uuid"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	if loginData.UserLog.Email == "" {
		redirectURL := "/login"
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // setting the max memory for a photo to 20 Mb
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
		beauty, _ := strconv.Atoi(r.FormValue("beauty-rating"))
		danger, _ := strconv.Atoi(r.FormValue("danger-rating"))
		err = os.MkdirAll("./assets/storage", os.ModePerm)
		if err != nil {
			http.Error(w, "Unable to create storage directory", http.StatusInternalServerError)
			return
		}
		files := r.MultipartForm.File["photos"]
		var imageUUIDs []string

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Unable to open file", http.StatusInternalServerError)
				return
			}
			defer file.Close()
			imageUUID, err := uuid.NewV4()
			if err != nil {
				http.Error(w, "Unable to generate UUID", http.StatusInternalServerError)
				return
			}
			imagePath := filepath.Join("./assets/storage", imageUUID.String()+filepath.Ext(fileHeader.Filename))
			newFile, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Unable to create file", http.StatusInternalServerError)
				return
			}
			defer newFile.Close()
			_, err = io.Copy(newFile, file)
			if err != nil {
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			imageUUIDs = append(imageUUIDs, imageUUID.String()+filepath.Ext(fileHeader.Filename))
		}
		categoryIds := r.Form["categories"]
		var categories []dbmanagement.Categorie
		for _, categoryId := range categoryIds {
			id, _ := strconv.Atoi(categoryId)
			categoryIn, _ := dbmanagement.DB.GetCategorie(id)
			categories = append(categories, categoryIn)
		}
		post, _ := loginData.UserLog.AddPost(loginData.UserLog.Email, loginData.UserLog.Password, r.FormValue("title"), r.FormValue("description"), imageUUIDs, danger, beauty, []int{})

		for _, category := range categories {
			post.AddToCategorie(category)
		}
		loginData.UserLog.AddNotification("Votre post "+r.FormValue("title")+" a bien été ajouté !", "success")
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

	loginData.UserLog, _, _ = dbmanagement.DB.ConnectToAccount(loginData.UserLog.Email, loginData.UserLog.Password)
	loginData.CategoriesAvailable = dbmanagement.DB.GetCategories(loginData.UserLog.Email, loginData.UserLog.Password)
	tmpl, err := template.ParseFiles("./assets/pages/addpost.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
