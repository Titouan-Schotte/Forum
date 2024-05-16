package dbmanagement

import (
	"fmt"
	"log"
)

func (db *DBForum) AddCategorie(email string, password string, nomCat string) (Categorie, bool) {
	// Vérifier les autorisations de l'utilisateur
	user, ok, _ := db.GetUser(email)
	fmt.Println(user.Password)
	if !ok || user.Password != password {
		return Categorie{}, false
	}

	// Vérifier si l'utilisateur est un administrateur
	if !user.IsAdmin {
		return Categorie{}, false
	}

	// Préparer la requête d'insertion de la nouvelle catégorie
	stmt, err := db.core.Prepare("INSERT INTO Categorie(Nom) VALUES(?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion de la catégorie:", err)
	}
	defer stmt.Close()

	// Exécuter la requête d'insertion de la nouvelle catégorie
	result, err := stmt.Exec(nomCat)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête d'insertion de la catégorie:", err)
	}

	// Obtenir l'ID de la nouvelle catégorie insérée
	catId, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Erreur lors de l'obtention de l'ID de la nouvelle catégorie insérée:", err)
	}

	// Créer une nouvelle structure Categorie avec les données fournies
	newCat := Categorie{
		Id:    int(catId),
		Nom:   nomCat,
		Posts: []Post{}, // Initialize posts as empty slice
	}

	// Retourner la nouvelle catégorie et true pour indiquer que l'opération a réussi
	return newCat, true
}

func (cat *Categorie) EditCategorie(newNameCat string) (Categorie, bool) {

	// Vérifier si l'utilisateur est un administrateur
	if !cat.UserManipuling.IsUserAdminGranted() {
		return Categorie{}, false
	}

	// Préparer la requête de mise à jour du nom de la catégorie
	stmt, err := DB.core.Prepare("UPDATE Categorie SET Nom = ? WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de modification de la catégorie:", err)
	}
	defer stmt.Close()

	// Exécuter la requête de mise à jour du nom de la catégorie
	_, err = stmt.Exec(newNameCat, cat.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de modification de la catégorie:", err)
	}

	// Créer une nouvelle structure Categorie avec les données fournies
	editedCat := Categorie{
		Id:    cat.Id,
		Nom:   newNameCat,
		Posts: nil, // Initialize posts as empty slice
	}

	// Retourner la catégorie modifiée et true pour indiquer que l'opération a réussi
	return editedCat, true
}

func (cat *Categorie) DeleteCategorie() bool {

	// Vérifier si l'utilisateur est un administrateur
	fmt.Println(cat.UserManipuling.IsUserAdminGranted())
	if !cat.UserManipuling.IsUserAdminGranted() {
		return false
	}

	// Préparer la requête de suppression de la catégorie
	stmt, err := DB.core.Prepare("DELETE FROM Categorie WHERE Id = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression de la catégorie:", err)
	}
	defer stmt.Close()

	// Exécuter la requête de suppression de la catégorie
	_, err = stmt.Exec(cat.Id)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression de la catégorie:", err)
	}

	// Retourner true pour indiquer que l'opération a réussi
	return true
}

func (db *DBForum) GetCategorie(email string, password string, id int) (Categorie, bool) {
	user, success := IsUserConnected(email, password)
	if !success {
		return Categorie{}, false
	}
	row := db.core.QueryRow("SELECT Nom FROM Categorie WHERE Id = ?", id)

	// Créer une structure User pour stocker les données de l'utilisateur
	var categorie Categorie
	// Scanner les valeurs des colonnes dans la structure User
	row.Scan(&categorie.Nom)

	categorie.Id = id
	categorie.UserManipuling = user

	categorie.Posts = DB.GetPostsOfCategory(categorie)

	return categorie, true
}

func (db *DBForum) GetCategories(email string, password string) []Categorie {
	// Créer une slice pour stocker les utilisateurs récupérés
	var categories []Categorie
	user, success := IsUserConnected(email, password)
	if !success {
		return []Categorie{}
	}
	rows, err := db.core.Query("SELECT Id,Nom FROM Categorie")
	if err != nil {
		return categories
	}
	defer rows.Close()

	// Parcourir les lignes résultantes
	for rows.Next() {
		var categorie Categorie
		// Scanner les valeurs des colonnes dans la structure User
		err := rows.Scan(&categorie.Id, &categorie.Nom)
		if err != nil {
			return categories
		}
		categorie.UserManipuling = user
		categorie.Posts = DB.GetPostsOfCategory(categorie)
		// Ajouter l'utilisateur à la slice
		categories = append(categories, categorie)
	}
	// Vérifier les erreurs éventuelles lors du parcours des lignes
	if err := rows.Err(); err != nil {
		return categories
	}

	return categories
}
