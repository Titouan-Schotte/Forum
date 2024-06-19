/*
Titouan Schotté

Utilities for db management
*/

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
	rows, err := db.core.Query("SELECT Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan FROM User")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Pseudo, &user.Email, &user.Password, &user.IsCertified, &user.IsModo, &user.IsAdmin, &user.IsBan)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (db *DBForum) GetUser(email string) (User, bool, string) {
	row := db.core.QueryRow("SELECT Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan FROM User WHERE Email=?", email)

	var user User
	err := row.Scan(&user.Pseudo, &user.Email, &user.Password, &user.IsCertified, &user.IsModo, &user.IsAdmin, &user.IsBan)
	if err != nil {
		// Si l'utilisateur n'est pas trouvé, retourner une structure User vide et false
		return User{}, false, "Utilisateur introuvable"
	}

	return user, true, ""
}

func GetUserBasicInfo(email string) User {
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
