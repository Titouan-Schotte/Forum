/*
Titouan Schotté

Utilities for commentaries management database SQLITE
*/

package dbmanagement

import (
	"log"
)

func (post *Post) AddComment(email string, password string, content string) (Comment, bool) {
	user, ok, _ := DB.GetUser(email)
	if !ok || user.Password != password || user.IsBan {
		return Comment{}, false
	}

	stmt, err := DB.core.Prepare("INSERT INTO Comment(Content, LikeCount, DislikeCount, PostId, AuthorEmail) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion du commentaire:", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(content, 0, 0, post.Id, user.Email)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête d'insertion du commentaire:", err)
	}

	commentId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Erreur lors de l'obtention de l'ID du nouveau commentaire inséré:", err)
	}

	newComment := Comment{
		Id:      int(commentId),
		Content: content,
		Author:  user,
		Like:    0,
		Dislike: 0,
	}

	return newComment, true
}

func DeleteComment(email string, password string, comment Comment) bool {
	user, ok, _ := DB.GetUser(email)
	if !ok || user.Password != password || user.IsBan {
		return false
	}

	if comment.Author.Email != user.Email {
		return false
	}

	stmt, err := DB.core.Prepare("DELETE FROM Comment WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression du commentaire:", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression du commentaire:", err)
	}

	return true
}

func (comment *Comment) EditComment(email string, password string) bool {
	user, ok, _ := DB.GetUser(email)
	if !ok || user.Password != password || user.IsBan {
		return false
	}

	stmt, err := DB.core.Prepare("UPDATE Comment SET Content = ?, LikeCount = ?, DislikeCount = ?, AuthorEmail = ?, PostId = ? WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de mise à jour du nombre de likes du commentaire:", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(comment.Content, comment.Like, comment.Dislike, comment.Author.Email, comment.PostOrigin.Id, comment.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de mise à jour du nombre de likes du commentaire:", err)
	}
	return true
}

func (post *Post) LoadComments() []Comment {
	rows, err := DB.core.Query("SELECT Id, Content, AuthorEmail, LikeCount, DislikeCount FROM Comment WHERE PostId = ?", post.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var comment Comment

		err := rows.Scan(&comment.Id, &comment.Content, &comment.Author.Email, &comment.Like, &comment.Dislike)
		if err != nil {
			return nil
		}
		comment.PostOrigin = *post
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	return comments
}

func (comment *Comment) GetAllLikesUsers() []User {

	rows, err := DB.core.Query("SELECT AuthorEmail FROM LikesComments WHERE CommentId = ?", comment.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		rows.Scan(&user.Email)
		user = GetUserBasicInfo(user.Email)
		users = append(users, user)
	}

	return users
}

func (comment *Comment) GetAllDislikesUsers() []User {
	rows, err := DB.core.Query("SELECT AuthorEmail FROM DislikesComments WHERE CommentId = ?", comment.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		rows.Scan(&user.Email)
		user = GetUserBasicInfo(user.Email)
		users = append(users, user)
	}

	return users
}

func (db *DBForum) GetCommentById(email, password string, id int) Comment {
	rows, err := db.core.Query("SELECT Id, Content, LikeCount, DislikeCount, AuthorEmail, PostId FROM Comment WHERE Id = ?", id)
	if err != nil {
		return Comment{}
	}
	defer rows.Close()

	var comment Comment

	if rows.Next() {
		var authorEmail string
		var postId int
		err := rows.Scan(&comment.Id, &comment.Content, &comment.Like, &comment.Dislike, &authorEmail, &postId)
		if err != nil {
			return Comment{}
		}

		comment.Author = GetUserBasicInfo(authorEmail)
		comment.PostOrigin = db.GetPostById(email, password, postId)
	}

	if err := rows.Err(); err != nil {
		return Comment{}
	}

	return comment
}
