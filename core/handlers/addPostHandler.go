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
	// Is Not logged in => redirect to login page
	if loginData.UserLog.Email == "" {
		redirectURL := "/login"
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	// Submitting formular
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // maxMemory 10MB
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		beauty, _ := strconv.Atoi(r.FormValue("beauty-rating"))
		danger, _ := strconv.Atoi(r.FormValue("danger-rating"))

		// Create the storage directory if it doesn't exist
		err = os.MkdirAll("./assets/storage", os.ModePerm)
		if err != nil {
			http.Error(w, "Unable to create storage directory", http.StatusInternalServerError)
			return
		}

		// Retrieve the files from the form
		files := r.MultipartForm.File["photos"]
		var imageUUIDs []string

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, "Unable to open file", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			// Generate a new UUID
			imageUUID, err := uuid.NewV4()
			if err != nil {
				http.Error(w, "Unable to generate UUID", http.StatusInternalServerError)
				return
			}
			imagePath := filepath.Join("./assets/storage", imageUUID.String()+filepath.Ext(fileHeader.Filename))

			// Create the new file
			newFile, err := os.Create(imagePath)
			if err != nil {
				http.Error(w, "Unable to create file", http.StatusInternalServerError)
				return
			}
			defer newFile.Close()

			// Copy the file contents to the new file
			_, err = io.Copy(newFile, file)
			if err != nil {
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}

			// Append the UUID to the list of image UUIDs
			imageUUIDs = append(imageUUIDs, imageUUID.String()+filepath.Ext(fileHeader.Filename))
		}

		loginData.UserLog.AddPost(loginData.UserLog.Email, loginData.UserLog.Password, r.FormValue("title"), r.FormValue("description"), imageUUIDs, danger, beauty, []int{1, 2})
	}

	loginData.UserLog, _, _ = dbmanagement.DB.ConnectToAccount(loginData.UserLog.Email, loginData.UserLog.Password)
	// Load the home page template
	tmpl, err := template.ParseFiles("./assets/pages/addpost.html")
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
