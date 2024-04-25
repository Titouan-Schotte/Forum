package main

import (
	"database/sql"
	"log"
)
import _ "github.com/mattn/go-sqlite3"

var dbname = "./brainstorm/forum.db"

type DBForum struct {
	core *sql.DB
}

func LoadDb() DBForum {
	dbIn, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base de données:", err)
	}
	return DBForum{core: dbIn}
}
func (db *DBForum) GetUsers() ([]User, error) {
	// Préparer la requête de sélection
	rows, err := db.core.Query("SELECT Id, Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan FROM User")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Créer une slice pour stocker les utilisateurs récupérés
	var users []User

	// Parcourir les lignes résultantes
	for rows.Next() {
		var user User
		// Scanner les valeurs des colonnes dans la structure User
		err := rows.Scan(&user.Id, &user.Pseudo, &user.Email, &user.Password, &user.IsCertified, &user.IsModo, &user.IsAdmin, &user.IsBan)
		if err != nil {
			return nil, err
		}
		// Ajouter l'utilisateur à la slice
		users = append(users, user)
	}

	// Vérifier les erreurs éventuelles lors du parcours des lignes
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DBForum) CreateAccount(pseudo string, email string, password string) User {
	// Préparer la requête d'insertion
	stmt, err := db.core.Prepare("INSERT INTO User(Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête d'insertion:", err)
	}
	defer stmt.Close()

	// Exécuter la requête d'insertion
	result, err := stmt.Exec(pseudo, email, password)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête d'insertion:", err)
	}
	// Obtenir l'ID du nouvel utilisateur inséré
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal("Erreur lors de l'obtention de l'ID du nouvel utilisateur:", err)
	}

	// Créer un nouvel utilisateur avec les données fournies et l'ID généré
	newUser := User{
		Id:          int(id),
		Pseudo:      pseudo,
		Email:       email,
		Password:    password,
		IsCertified: false,
		IsModo:      false,
		IsAdmin:     false,
		IsBan:       false,
	}

	return newUser
}

func (db *DBForum) ConnectToAccount(email string, password string) (User, bool) {
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
