/*
Titouan Schotté

Utilities for posts management database SQLITE
*/

package dbmanagement

import (
	"log"
	"strings"
	"time"
)

func (user *User) AddPost(email string, password string, titlePost string, descriptionPost string, photosPost []string, dangerPost int, beauty int, categorie []int) (Post, bool) {
	// Vérifier les autorisations de l'utilisateur
	if user.Email != email || user.Password != password || user.IsBan {
		return Post{}, false
	}
	// Préparer la requête d'insertion du nouveau commentaire
	stmt, err := DB.core.Prepare("INSERT INTO Post(Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, DatePost) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion du commentaire:", err)
	}
	defer stmt.Close()

	photoText := strings.Join(photosPost, ";")

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	result, err := stmt.Exec(titlePost, descriptionPost, dangerPost, beauty, 0, 0, user.Email, photoText, formattedTime)

	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête d'insertion du commentaire:", err)
	}

	postId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Erreur lors de l'obtention de l'ID du nouveau commentaire inséré:", err)
	}

	newPost := Post{
		Id:          int(postId),
		Title:       titlePost,
		Description: descriptionPost,
		Photos:      photosPost,
		Danger:      dangerPost,
		Beauty:      beauty,
		Author:      *user,
		Date:        formattedTime,
	}
	for _, v := range categorie {
		cat, _ := DB.GetCategorie(v)
		newPost.AddToCategorie(cat)
	}
	newPost.Categories = newPost.GetCategories()
	return newPost, true
}

func (post *Post) EditPost(email string, password string) bool {
	user, success := IsUserConnected(email, password)
	if !success || user.IsBan || post.Author.Email != user.Email {
		return false
	}
	photoText := strings.Join(post.Photos, ";")
	stmt, err := DB.core.Prepare("UPDATE Post SET Title = ?, Description = ?, Danger = ?, Beauty = ?, LikeCount = ?, DislikeCount = ?, AuthorEmail = ?, Photos = ?, DatePost = ? WHERE Id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Title, post.Description, post.Danger, post.Beauty, post.Like, post.Dislike, post.Author.Email, photoText, post.Date, post.Id)
	if err != nil {
		return false
	}
	return true
}

func (post *Post) LikePost(email string, password string) bool {
	post.Like++
	post.EditPost(email, password)

	user, success := IsUserConnected(email, password)
	if !success || user.IsBan || post.Author.Email != user.Email {
		return false
	}
	var count int
	DB.core.QueryRow("SELECT COUNT(*) FROM Likes WHERE AuthorEmail = ? AND PostId = ?", email, post.Id).Scan(&count)
	if count == 0 {
		_, err := DB.core.Exec("INSERT INTO Likes(PostId, AuthorEmail) VALUES(?, ?)", post.Id, user.Email)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (post *Post) DislikePost(email string, password string) bool {
	post.Dislike++
	post.EditPost(email, password)

	user, success := IsUserConnected(email, password)
	if !success || user.IsBan || post.Author.Email != user.Email {
		return false
	}
	var count int
	DB.core.QueryRow("SELECT COUNT(*) FROM Dislikes WHERE AuthorEmail = ? AND PostId = ?", email, post.Id).Scan(&count)
	if count == 0 {
		_, err := DB.core.Exec("INSERT INTO Dislikes(PostId, AuthorEmail) VALUES(?, ?)", post.Id, user.Email)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func (post *Post) DeletePost(email string) bool {
	if post.Author.Email != email || post.Author.IsBan {
		return false
	}
	stmt, err := DB.core.Prepare("DELETE FROM Post WHERE Id=?")
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.Id)
	if err != nil {
		return false
	}
	return true
}

func (userModo *User) DeletePostModo(post Post) bool {
	moderator, ok := IsUserConnected(userModo.Email, userModo.Password)
	if !ok || !moderator.IsModo {
		// Si le compte n'est pas un modérateur, retourner false
		return false
	}
	if !post.Author.IsModo || !post.Author.IsAdmin {
		stmt, err := DB.core.Prepare("DELETE FROM Post WHERE Id=?")
		if err != nil {
			log.Fatal("Erreur lors de la préparation de la requête de suppression:", err)
			return false
		}
		defer stmt.Close()

		_, err = stmt.Exec(post.Id)
		if err != nil {
			log.Fatal("Erreur lors de l'exécution de la requête de suppression:", err)
			return false
		}
	}
	return false
}

func (db *DBForum) GetPostsOfCategory(categorie Categorie) []Post {
	rows, err := db.core.Query("SELECT PostId FROM PostCategorie WHERE CategorieId = ?", categorie.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		var photos string
		rows.Scan(&post.Id)
		rowsPost, err := db.core.Query("SELECT Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos,DatePost FROM Post WHERE PostId = ?", post.Id)
		if err != nil {
			return nil
		}
		defer rowsPost.Close()
		err2 := rowsPost.Scan(&post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err2 != nil {
			return nil
		}

		post.Author = GetUserBasicInfo(post.AuthorEmail)
		post.Photos = strings.Split(photos, ";")

		post.Comments = post.LoadComments()
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil
	}

	return posts
}

func (post *Post) GetAllLikesUsers() []User {
	rows, err := DB.core.Query("SELECT AuthorEmail FROM Likes WHERE PostId = ?", post.Id)
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

func (post *Post) GetAllDislikesUsers() []User {
	rows, err := DB.core.Query("SELECT AuthorEmail FROM Dislikes WHERE PostId = ?", post.Id)
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

func (db *DBForum) GetPostById(email, password string, id int) Post {
	rows, err := db.core.Query("SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, DatePost FROM Post WHERE Id = ?", id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de récupération du post:", err)
		return Post{}
	}
	defer rows.Close()
	var post Post
	if rows.Next() {
		var photos string
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			log.Fatal("Erreur lors du scan des données du post:", err)
			return Post{}
		}
		post.Author = GetUserBasicInfo(post.AuthorEmail)
		post.Photos = strings.Split(photos, ";")
		post.Comments = post.LoadComments()
		post.Categories = post.GetCategories()
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Erreur lors du parcours des résultats:", err)
		return Post{}
	}
	return post
}

func (db *DBForum) GetMostRecentsPosts(numberOfPost int) []Post {
	rows, err := db.core.Query("SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, DatePost FROM Post ORDER BY DatePost DESC LIMIT ?", numberOfPost)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		var photos string
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			return nil
		}
		post.Author = GetUserBasicInfo(post.AuthorEmail)
		post.Photos = strings.Split(photos, ";")
		post.Categories = post.GetCategories()
		post.Comments = post.LoadComments()
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil
	}

	return posts
}

func (db *DBForum) GetTopPosts(numberOfPost int) []Post {
	rows, err := db.core.Query("SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, DatePost FROM Post ORDER BY LikeCount DESC LIMIT ?", numberOfPost)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		var photos string
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			return nil
		}
		post.Author = GetUserBasicInfo(post.AuthorEmail)
		post.Photos = strings.Split(photos, ";")
		post.Categories = post.GetCategories()
		post.Comments = post.LoadComments()
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil
	}

	return posts
}
func (db *DBForum) GetRandomPosts(numberOfPost int) []Post {
	query := "SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, DatePost FROM Post ORDER BY RANDOM() LIMIT ?"
	rows, err := db.core.Query(query, numberOfPost)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		var photos string
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			return nil
		}
		post.Author = GetUserBasicInfo(post.AuthorEmail)
		post.Photos = strings.Split(photos, ";")
		post.Categories = post.GetCategories()
		post.Comments = post.LoadComments()
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil
	}

	return posts
}
