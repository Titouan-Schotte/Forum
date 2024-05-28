package dbmanagement

import (
	"log"
	"strings"
	"time"
)

func (user *User) AddPost(email string, password string, titlePost string, descriptionPost string, photosPost []string, dangerPost int, beauty int, categorie Categorie) (Post, bool) {
	// Vérifier les autorisations de l'utilisateur
	if user.Email != email || user.Password != password || user.IsBan {
		return Post{}, false
	}
	// Préparer la requête d'insertion du nouveau commentaire
	stmt, err := DB.core.Prepare("INSERT INTO Post(Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, Categorie, DatePost) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion du commentaire:", err)
	}
	defer stmt.Close()

	photoText := strings.Join(photosPost, ";")

	// Obtenir la date et l'heure actuelles
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	// Exécuter la requête d'insertion du nouveau commentaire
	result, err := stmt.Exec(titlePost, descriptionPost, dangerPost, beauty, 0, 0, user.Email, photoText, categorie.Id, formattedTime)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête d'insertion du commentaire:", err)
	}

	// Obtenir l'ID du nouveau commentaire inséré
	postId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Erreur lors de l'obtention de l'ID du nouveau commentaire inséré:", err)
	}

	// Créer un nouveau post avec les données fournies
	newPost := Post{
		Id:          int(postId),
		Title:       titlePost,
		Description: descriptionPost,
		Photos:      photosPost,
		Danger:      dangerPost,
		Beauty:      beauty,
		Author:      *user,
		Categorie:   categorie,
		Date:        formattedTime,
	}

	// Retourner le nouveau post et true pour indiquer que l'opération a réussi
	return newPost, true
}

func (post *Post) EditPost(email string, password string) bool {
	// Vérifier les autorisations de l'utilisateur
	user, success := IsUserConnected(email, password)
	if !success || user.IsBan || post.Author.Email != user.Email {
		return false
	}
	photoText := strings.Join(post.Photos, ";")
	// Préparer la requête de mise à jour du nombre de likes du commentaire
	stmt, err := DB.core.Prepare("UPDATE Post SET Title = ?, Description = ?, Danger = ?, Beauty = ?, LikeCount = ?, DislikeCount = ?, AuthorEmail = ?, Photos = ?, Categorie = ?, DatePost = ? WHERE Id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()

	// Exécuter la requête de mise à jour du nombre de likes du commentaire
	_, err = stmt.Exec(post.Title, post.Description, post.Danger, post.Beauty, post.Like, post.Dislike, post.Author.Email, photoText, post.Categorie.Id, post.Date, post.Id)
	if err != nil {
		return false
	}
	// Retourner true pour indiquer que l'édition du post a réussi
	return true
}

func (post *Post) LikePost(email string, password string) bool {

	// Incrémenter le compteur de likes du post
	post.Like++
	post.EditPost(email, password)

	user, success := IsUserConnected(email, password)
	if !success || user.IsBan || post.Author.Email != user.Email {
		return false
	}
	//Enregistrer le like
	rows, _ := DB.core.Query("SELECT PostId, AuthorEmail FROM Likes WHERE AuthorEmail = ? AND PostId = ?", email, post.Id)

	if !rows.Next() && !rows.Next() {
		stmt, _ := DB.core.Prepare("INSERT INTO Likes(PostId, AuthorEmail) VALUES(?, ?)")
		defer stmt.Close()
		stmt.Exec(post.Id, user.Email)
		return true

	}
	// Retourner le post mis à jour et true pour indiquer que l'opération a réussi
	return false
}

func (post *Post) DislikePost(email string, password string) bool {
	// Incrémenter le compteur de dislikes du post
	post.Dislike++
	post.EditPost(email, password)
	user, success := IsUserConnected(email, password)
	if !success || user.IsBan || post.Author.Email != user.Email {
		return false
	}
	//Enregistrer le like
	rows, _ := DB.core.Query("SELECT PostId, AuthorEmail FROM Dislikes WHERE AuthorEmail = ? AND PostId = ?", user.Email, post.Id)
	if !rows.Next() && !rows.Next() {
		stmt, _ := DB.core.Prepare("INSERT INTO Dislikes(PostId, AuthorEmail) VALUES(?, ?)")
		defer stmt.Close()
		stmt.Exec(post.Id, user.Email)
		return true

	}
	// Retourner le post mis à jour et true pour indiquer que l'opération a réussi
	return false
}

func (post *Post) DeletePost(email string) bool {
	// Vérifier les autorisations de l'utilisateur
	if post.Author.Email != email || post.Author.IsBan {
		return false
	}
	stmt, err := DB.core.Prepare("DELETE FROM Post WHERE Id=?")
	if err != nil {
		return false
	}
	defer stmt.Close()

	// Exécuter la requête de suppression
	_, err = stmt.Exec(post.Id)
	if err != nil {
		return false
	}
	// Retourner false si le post n'a pas été trouvé
	return true
}

func (userModo *User) DeletePostModo(post Post) bool {
	// Vérifier si le compte qui initie l'action est un modérateur
	moderator, ok := IsUserConnected(userModo.Email, userModo.Password)
	if !ok || !moderator.IsModo {
		// Si le compte n'est pas un modérateur, retourner false
		return false
	}

	// Supprimer le post uniquement si l'auteur du post n'est pas un modérateur
	if !post.Author.IsModo || !post.Author.IsAdmin {
		stmt, err := DB.core.Prepare("DELETE FROM Post WHERE Id=?")
		if err != nil {
			log.Fatal("Erreur lors de la préparation de la requête de suppression:", err)
			return false
		}
		defer stmt.Close()

		// Exécuter la requête de suppression
		_, err = stmt.Exec(post.Id)
		if err != nil {
			log.Fatal("Erreur lors de l'exécution de la requête de suppression:", err)
			return false
		}
	}

	// Retourner false si l'auteur du post est un modérateur
	return false
}

func (db *DBForum) GetPostsOfCategory(categorie Categorie) []Post {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := db.core.Query("SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos,DatePost FROM Post WHERE Categorie = ?", categorie.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var posts []Post

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var post Post
		var photos string // Stockage des photos en tant que chaîne séparée par des points-virgules

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			return nil
		}

		post.Author = GetUserBasicInfo(post.AuthorEmail)

		// Diviser la chaîne de photos en une slice de chaînes
		post.Photos = strings.Split(photos, ";")

		post.Comments = post.LoadComments()
		post.Categorie = categorie
		// Ajout du post à la slice des posts
		posts = append(posts, post)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	if err := rows.Err(); err != nil {
		return nil
	}

	return posts
}

func (post *Post) GetAllLikesUsers() []User {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT AuthorEmail FROM Likes WHERE PostId = ?", post.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var users []User

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var user User

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&user.Email)
		user = GetUserBasicInfo(user.Email)
		// Ajout du post à la slice des posts
		users = append(users, user)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return users
}

func (post *Post) GetAllDislikesUsers() []User {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT AuthorEmail FROM Dislikes WHERE PostId = ?", post.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var users []User

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var user User

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&user.Email)
		user = GetUserBasicInfo(user.Email)
		// Ajout du post à la slice des posts
		users = append(users, user)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return users
}

func (db *DBForum) GetPostById(email, password string, id int) Post {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := db.core.Query("SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, Categorie, DatePost FROM Post WHERE Id = ?", id)
	if err != nil {
		return Post{}
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var post Post

	// Parcours des résultats et création des structures Post
	if rows.Next() {
		var photos string // Stockage des photos en tant que chaîne séparée par des points-virgules
		var catInt int
		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &catInt, &post.Date)
		if err != nil {
			return Post{}
		}

		post.Author = GetUserBasicInfo(post.AuthorEmail)

		// Diviser la chaîne de photos en une slice de chaînes
		post.Photos = strings.Split(photos, ";")

		post.Comments = post.LoadComments()
		post.Categorie, _ = db.GetCategorie(email, password, id)
		// Ajout du post à la slice des posts
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	if err := rows.Err(); err != nil {
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

	// Création d'une slice pour stocker les posts récupérés
	var posts []Post

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var post Post
		var photos string // Stockage des photos en tant que chaîne séparée par des points-virgules

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			return nil
		}

		post.Author = GetUserBasicInfo(post.AuthorEmail)

		// Diviser la chaîne de photos en une slice de chaînes
		post.Photos = strings.Split(photos, ";")

		post.Comments = post.LoadComments()
		// Ajout du post à la slice des posts
		posts = append(posts, post)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
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

	// Création d'une slice pour stocker les posts récupérés
	var posts []Post

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var post Post
		var photos string // Stockage des photos en tant que chaîne séparée par des points-virgules

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			return nil
		}

		post.Author = GetUserBasicInfo(post.AuthorEmail)

		// Diviser la chaîne de photos en une slice de chaînes
		post.Photos = strings.Split(photos, ";")

		post.Comments = post.LoadComments()
		// Ajout du post à la slice des posts
		posts = append(posts, post)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	if err := rows.Err(); err != nil {
		return nil
	}

	return posts
}
func (db *DBForum) GetRandomPosts(numberOfPost int) []Post {
	// Exécution de la requête SQL pour récupérer des posts aléatoires et limiter le nombre de résultats
	query := "SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, DatePost FROM Post ORDER BY RANDOM() LIMIT ?"
	rows, err := db.core.Query(query, numberOfPost)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var posts []Post

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var post Post
		var photos string // Stockage des photos en tant que chaîne séparée par des points-virgules

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos, &post.Date)
		if err != nil {
			return nil
		}

		post.Author = GetUserBasicInfo(post.AuthorEmail)

		// Diviser la chaîne de photos en une slice de chaînes
		post.Photos = strings.Split(photos, ";")

		post.Comments = post.LoadComments()
		// Ajout du post à la slice des posts
		posts = append(posts, post)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	if err := rows.Err(); err != nil {
		return nil
	}

	return posts
}
