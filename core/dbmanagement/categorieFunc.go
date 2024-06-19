/*
Titouan Schotté

Utilities for categories management database SQLITE
*/

package dbmanagement

import (
	"log"
)

func (db *DBForum) AddCategorie(email string, password string, nomCat string) (Categorie, bool) {
	user, ok, _ := db.GetUser(email)
	if !ok || user.Password != password {
		return Categorie{}, false
	}

	if !user.IsAdmin {
		return Categorie{}, false
	}

	stmt, err := db.core.Prepare("INSERT INTO Categorie(Nom) VALUES(?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion de la catégorie:", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(nomCat)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête d'insertion de la catégorie:", err)
	}

	catId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Erreur lors de l'obtention de l'ID de la nouvelle catégorie insérée:", err)
	}

	newCat := Categorie{
		Id:    int(catId),
		Nom:   nomCat,
		Posts: []Post{},
	}

	return newCat, true
}

func (cat *Categorie) EditCategorie(email, password, newNameCat string) (Categorie, bool) {

	user, _, _ := DB.ConnectToAccount(email, password)
	if !user.IsUserAdminGranted() {
		return Categorie{}, false
	}

	stmt, err := DB.core.Prepare("UPDATE Categorie SET Nom = ? WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de modification de la catégorie:", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(newNameCat, cat.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de modification de la catégorie:", err)
	}

	editedCat := Categorie{
		Id:    cat.Id,
		Nom:   newNameCat,
		Posts: nil,
	}

	return editedCat, true
}

func (cat *Categorie) DeleteCategorie(email, password string) bool {

	user, _, _ := DB.ConnectToAccount(email, password)
	if !user.IsUserAdminGranted() {
		return false
	}

	stmt, err := DB.core.Prepare("DELETE FROM Categorie WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression de la catégorie:", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(cat.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression de la catégorie:", err)
	}

	return true
}

func (db *DBForum) GetCategorie(id int) (Categorie, bool) {

	row := db.core.QueryRow("SELECT Nom FROM Categorie WHERE Id = ?", id)

	var categorie Categorie
	row.Scan(&categorie.Nom)

	categorie.Id = id

	categorie.Posts = DB.GetPostsOfCategory(categorie)

	return categorie, true
}

func (db *DBForum) GetCategories(email string, password string) []Categorie {
	var categories []Categorie

	rows, err := db.core.Query("SELECT Id,Nom FROM Categorie")
	if err != nil {
		return categories
	}
	defer rows.Close()

	for rows.Next() {
		var categorie Categorie
		err := rows.Scan(&categorie.Id, &categorie.Nom)
		if err != nil {
			return categories
		}
		categorie.Posts = DB.GetPostsOfCategory(categorie)
		categories = append(categories, categorie)
	}
	if err := rows.Err(); err != nil {
		return categories
	}

	return categories
}

func (post *Post) AddToCategorie(categorie Categorie) {
	stmt, err := DB.core.Prepare("INSERT INTO PostCategorie(PostId, CategorieId) VALUES(?, ?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion de la catégorie:", err)
	}
	defer stmt.Close()

	stmt.Exec(post.Id, categorie.Id)

}
func (post *Post) DeleteOfCategorie(categorie Categorie) {
	stmt, err := DB.core.Prepare("DELETE FROM PostCategorie WHERE PostId = ? AND CategorieId = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion de la catégorie:", err)
	}
	defer stmt.Close()

	stmt.Exec(post.Id, categorie.Id)
}

func (post *Post) GetCategories() []Categorie {
	rows, _ := DB.core.Query("SELECT CategorieId FROM PostCategorie WHERE PostId = ?", post.Id)
	defer rows.Close()
	var categories []Categorie
	// Parcourir les lignes résultantes
	for rows.Next() {
		var categorie Categorie
		err := rows.Scan(&categorie.Id)
		if err != nil {
			return categories
		}
		categorie, _ = DB.GetCategorie(categorie.Id)
		categories = append(categories, categorie)
	}
	if err := rows.Err(); err != nil {
		return categories
	}

	return categories
}
