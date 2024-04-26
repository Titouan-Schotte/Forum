package dbmanagement

func (db *DBForum) AddPost(email string, password string, titlePost string, descriptionPost string, photosPost []string, dangerPost int, beauty int, categories []Categorie) (Post, bool) {
	return Post{}, true
}

func EditPost(email string, password string, post Post, titlePost string, descriptionPost string, photosPost []string, dangerPost int, beauty int, categories []Categorie) bool {
	return true
}

func LikePost(email string, password string, post Post) (Post, bool) {
	return Post{}, true
}

func DislikePost(email string, password string, post Post) (Post, bool) {
	return Post{}, true
}

func DeletePost(email string, password string, post Post) bool {
	return true
}

func DeletePostModo(emailModo string, passwordModo string, post Post) bool {
	return true
}
