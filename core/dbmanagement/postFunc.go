package dbmanagement

import (
	"log"
	"strings"
)

func (user *User) AddPost(email string, password string, titlePost string, descriptionPost string, photosPost []string, dangerPost int, beauty int, categorie Categorie) (Post, bool) {
	// Vérifier les autorisations de l'utilisateur
	if user.Email != email || user.Password != password || user.IsBan {
		return Post{}, false
	}
	// Préparer la requête d'insertion du nouveau commentaire
	stmt, err := DB.core.Prepare("INSERT INTO Post(Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos, Categorie) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion du commentaire:", err)
	}
	defer stmt.Close()

	photoText := strings.Join(photosPost, ";")
	// Exécuter la requête d'insertion du nouveau commentaire
	result, err := stmt.Exec(titlePost, descriptionPost, dangerPost, beauty, 0, 0, user.Email, photoText, categorie.Id)
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
	stmt, err := DB.core.Prepare("UPDATE Post SET Title = ?, Description = ?, Danger = ?, Beauty = ?, LikeCount = ?, DislikeCount = ?, AuthorEmail = ?, Photos = ?, Categorie = ? WHERE Id = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()

	// Exécuter la requête de mise à jour du nombre de likes du commentaire
	_, err = stmt.Exec(post.Title, post.Description, post.Danger, post.Beauty, post.Like, post.Dislike, post.Author.Email, photoText, post.Categorie.Id, post.Id)
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
	// Retourner le post mis à jour et true pour indiquer que l'opération a réussi
	return true
}

func (post *Post) DislikePost(email string, password string) bool {
	// Incrémenter le compteur de dislikes du post
	post.Dislike++
	post.EditPost(email, password)
	// Retourner le post mis à jour et true pour indiquer que l'opération a réussi
	return true
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
	rows, err := db.core.Query("SELECT Id, Title, Description, Danger, Beauty, LikeCount, DislikeCount, AuthorEmail, Photos FROM Post WHERE Categorie = ?", categorie.Id)
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
		err := rows.Scan(&post.Id, &post.Title, &post.Description, &post.Danger, &post.Beauty, &post.Like, &post.Dislike, &post.AuthorEmail, &photos)
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
