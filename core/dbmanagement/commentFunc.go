package dbmanagement

import (
	"fmt"
	"log"
)

func (post *Post) AddComment(email string, password string, content string) (Comment, bool) {
	// Vérifier les autorisations de l'utilisateur
	user, ok, _ := DB.GetUser(email)
	if !ok || user.Password != password || user.IsBan {
		return Comment{}, false
	}

	// Préparer la requête d'insertion du nouveau commentaire
	stmt, err := DB.core.Prepare("INSERT INTO Comment(Content, LikeCount, DislikeCount, PostId, AuthorEmail) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion du commentaire:", err)
	}
	defer stmt.Close()

	// Exécuter la requête d'insertion du nouveau commentaire
	result, err := stmt.Exec(content, 0, 0, post.Id, user.Email)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête d'insertion du commentaire:", err)
	}

	// Obtenir l'ID du nouveau commentaire inséré
	commentId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Erreur lors de l'obtention de l'ID du nouveau commentaire inséré:", err)
	}

	// Créer une nouvelle structure Comment avec les données fournies
	newComment := Comment{
		Id:      int(commentId),
		Content: content,
		Author:  user,
		Like:    0,
		Dislike: 0,
	}

	// Retourner le nouveau commentaire et true pour indiquer que l'opération a réussi
	return newComment, true
}

func DeleteComment(email string, password string, comment Comment) bool {
	// Vérifier les autorisations de l'utilisateur
	user, ok, _ := DB.GetUser(email)
	if !ok || user.Password != password || user.IsBan {
		return false
	}

	// Vérifier si l'utilisateur est l'auteur du commentaire
	if comment.Author.Email != user.Email {
		return false
	}

	// Préparer la requête de suppression du commentaire
	stmt, err := DB.core.Prepare("DELETE FROM Comment WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression du commentaire:", err)
	}
	defer stmt.Close()

	// Exécuter la requête de suppression du commentaire
	_, err = stmt.Exec(comment.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression du commentaire:", err)
	}

	// Retourner true pour indiquer que l'opération a réussi
	return true
}

func (comment *Comment) EditComment(email string, password string) bool {
	user, ok, _ := DB.GetUser(email)
	if !ok || user.Password != password || user.IsBan {
		return false
	}

	// Préparer la requête de mise à jour du nombre de likes du commentaire
	stmt, err := DB.core.Prepare("UPDATE Comment SET Content = ?, LikeCount = ?, DislikeCount = ?, AuthorEmail = ?, PostId = ? WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de mise à jour du nombre de likes du commentaire:", err)
	}
	defer stmt.Close()

	fmt.Println(comment.PostOrigin.Id)
	// Exécuter la requête de mise à jour du nombre de likes du commentaire
	_, err = stmt.Exec(comment.Content, comment.Like, comment.Dislike, comment.Author.Email, comment.PostOrigin.Id, comment.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de mise à jour du nombre de likes du commentaire:", err)
	}
	return true
}

func (comment *Comment) LikeComment(email string, password string) bool {
	// Incrémenter le nombre de likes dans la structure Comment
	comment.Like++
	comment.EditComment(email, password)
	// Retourner le commentaire mis à jour et true pour indiquer que l'opération a réussi
	return true
}
func (comment *Comment) DislikeComment(email string, password string) bool {
	// Incrémenter le nombre de likes dans la structure Comment
	comment.Dislike++
	comment.EditComment(email, password)
	// Retourner le commentaire mis à jour et true pour indiquer que l'opération a réussi
	return true
}

func (post *Post) LoadComments() []Comment {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les commentaires liés au post donné
	rows, err := DB.core.Query("SELECT Id, Content, AuthorEmail, LikeCount, DislikeCount FROM Comment WHERE PostId = ?", post.Id)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les commentaires récupérés
	var comments []Comment

	// Parcours des résultats et création des structures Comment
	for rows.Next() {
		var comment Comment

		// Scan des colonnes de la table Comment dans les champs correspondants de la structure Comment
		err := rows.Scan(&comment.Id, &comment.Content, &comment.Author.Email, &comment.Like, &comment.Dislike)
		if err != nil {
			return nil
		}
		comment.PostOrigin = *post
		// Ajouter le commentaire à la slice des commentaires
		comments = append(comments, comment)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	if err := rows.Err(); err != nil {
		return nil
	}

	return comments
}
