package dbmanagement

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
