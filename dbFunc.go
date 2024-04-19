package Forum

func CreateAccount(pseudo string, email string, password string) User {
	return User{}
}

func ConnectToAccount(email string, password string) (User, bool) {
	return User{}, true
}

func DeleteAccount(email string, password string) bool {
	return true
}

func DeleteAccountModo(emailModo string, passwordModo string, emailTarget string) bool {
	return true
}
func BanAccountModo(emailModo string, passwordModo string, emailTarget string) bool {
	return true
}

func AddCategorie(email string, password string, nomCat string) (Categorie, bool) {
	return Categorie{}, true
}

func EditCategorie(email string, password string, idCat int, newNameCat string) (Categorie, bool) {
	return Categorie{}, true
}

func DeleteCategorie(email string, password string, idCat int) bool {
	return true
}

func AddPost(email string, password string, titlePost string, descriptionPost string, photosPost []string, dangerPost int, beauty int, categories []Categorie) (Post, bool) {
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

func AddComment(email string, password string, post Post, content string) (Comment, bool) {
	return Comment{}, true
}

func DeleteComment(email string, password string, post Post, content string) bool {
	return true
}

func LikeComment(email string, password string, comment Comment) (Comment, bool) {
	return Comment{}, true
}

func DislikeComment(email string, password string, comment Comment) (Comment, bool) {
	return Comment{}, true
}
