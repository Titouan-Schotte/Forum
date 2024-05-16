package dbmanagement

import (
	"database/sql"
	"log"
)
import _ "github.com/mattn/go-sqlite3"

var dbname = "./assets/db/forum.db"

type DBForum struct {
	core *sql.DB
}

var DB = LoadDb()

func LoadDb() DBForum {
	dbIn, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal("Erreur lors de l'ouverture de la base de données:", err)
	}
	return DBForum{core: dbIn}
}
func (db *DBForum) GetUsers() ([]User, error) {
	// Préparer la requête de sélection
	rows, err := db.core.Query("SELECT Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan FROM User")
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
		err := rows.Scan(&user.Pseudo, &user.Email, &user.Password, &user.IsCertified, &user.IsModo, &user.IsAdmin, &user.IsBan)
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

func (db *DBForum) GetUser(email string) (User, bool, string) {
	// Préparer la requête de sélection pour récupérer l'utilisateur correspondant à l'email et au mot de passe
	row := db.core.QueryRow("SELECT Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan FROM User WHERE Email=?", email)

	// Créer une structure User pour stocker les données de l'utilisateur
	var user User
	// Scanner les valeurs des colonnes dans la structure User
	err := row.Scan(&user.Pseudo, &user.Email, &user.Password, &user.IsCertified, &user.IsModo, &user.IsAdmin, &user.IsBan)
	if err != nil {
		// Si l'utilisateur n'est pas trouvé, retourner une structure User vide et false
		return User{}, false, "Utilisateur introuvable"
	}

	// Si l'utilisateur est trouvé, retourner la structure User et true
	return user, true, ""
}

func GetUserBasicInfo(email string) User {
	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT Pseudo, Email, IsCertified, IsModo, IsAdmin, IsBan FROM User WHERE Email = ?", email)
	if err != nil {
		return User{}
	}
	defer rows.Close()
	if rows.Next() {
		var user User
		err := rows.Scan(&user.Pseudo, &user.Email, &user.IsCertified, &user.IsModo, &user.IsAdmin, &user.IsBan)
		if err != nil {
			return User{}
		}

		return user
	}
	return User{}
}
