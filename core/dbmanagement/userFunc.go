/*
Titouan Schotté

Utilities for users management database SQLITE
*/

package dbmanagement

import (
	"fmt"
	"log"
)

func (db *DBForum) CreateAccount(pseudo string, email string, password string) (User, bool, string) {
	stmt, err := db.core.Prepare("INSERT INTO User(Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print("Erreur lors de la préparation de la requête d'insertion:", err)
		return User{}, false, "Erreur de la base de données."
	}
	defer stmt.Close()
	_, err = stmt.Exec(pseudo, email, password, false, false, false, false)
	if err != nil {
		log.Print("Erreur lors de l'exécution de la requête d'insertion:", err)
		return User{}, false, "Erreur de la base de données."
	}
	newUser := User{
		Pseudo:      pseudo,
		Email:       email,
		Password:    password,
		IsCertified: false,
		IsModo:      false,
		IsAdmin:     false,
		IsBan:       false,
	}
	return newUser, true, ""
}

func (db *DBForum) ConnectToAccount(email string, password string) (User, bool, string) {
	row := db.core.QueryRow("SELECT Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan FROM User WHERE Email=? AND Password=?", email, password)
	var user User
	err := row.Scan(&user.Pseudo, &user.Email, &user.Password, &user.IsCertified, &user.IsModo, &user.IsAdmin, &user.IsBan)
	if err != nil {
		// Si l'utilisateur n'est pas trouvé, retourner une structure User vide et false
		return User{}, false, "Utilisateur introuvable / Mauvais mot de passe."
	}
	user.Followers = user.GetAllFollowers()
	user.Subscription = user.GetAllFollowedAccount()
	user.Posts = user.GetAllUserPosts()
	user.Likes = user.GetAllLikedPosts()
	user.Notifications = user.GetAllNotifications()
	user.TotalLikes = 0
	for _, v := range user.Posts {
		user.TotalLikes += v.Like
	}
	return user, true, ""
}
func (user *User) EditUser() bool {
	stmt, err := DB.core.Prepare("UPDATE User SET Pseudo = ?, Email = ?, Password = ?, IsCertified = ?, IsModo = ?, IsAdmin = ?, IsBan = ? WHERE Email = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Pseudo, user.Email, user.Password, user.IsCertified, user.IsModo, user.IsAdmin, user.IsBan, user.Email)
	if err != nil {
		return false
	}
	return true
}
func (user *User) DeleteAccount() (bool, string) {
	stmt, err := DB.core.Prepare("DELETE FROM User WHERE Email=? AND Password=?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression:", err)
		return false, "Erreur de la base de données."
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Email, user.Password)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression:", err)
		return false, "Erreur de la base de données."
	}
	return true, ""
}

func IsUserConnected(emailLocalstorage, passwordLocalstorage string) (User, bool) {
	user, ok, _ := DB.GetUser(emailLocalstorage)
	return user, ok && user.Password == passwordLocalstorage
}

func (user *User) IsUserAdminGranted() bool {
	return user.IsAdmin
}
func (user *User) IsUserModoGranted() bool {
	return user.IsModo
}

func (user *User) GetAllLikedPosts() []Post {
	rows, err := DB.core.Query("SELECT PostId FROM Likes WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		rows.Scan(&post.Id)
		post = DB.GetPostById(user.Email, user.Password, post.Id)
		posts = append(posts, post)
	}
	return posts
}

func (user *User) GetAllDislikedPosts() []Post {
	rows, err := DB.core.Query("SELECT PostId FROM Dislikes WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		rows.Scan(&post.Id)
		post = DB.GetPostById(user.Email, user.Password, post.Id)
		posts = append(posts, post)
	}
	return posts
}

func (user *User) GetAllUserPosts() []Post {
	rows, err := DB.core.Query("SELECT Id FROM Post WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var posts []Post
	for rows.Next() {
		var post Post
		rows.Scan(&post.Id)
		post = DB.GetPostById(user.Email, user.Password, post.Id)
		posts = append(posts, post)
	}
	return posts
}

func (user *User) Follow(emailTarget string) bool {
	tx, err := DB.core.Begin()
	if err != nil {
		fmt.Println("Erreur lors du début de la transaction:", err)
		return false
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	stmt, err := tx.Prepare("SELECT AuthorEmail, ReceiverEmail FROM Follow WHERE AuthorEmail = ? AND ReceiverEmail = ?")
	if err != nil {
		fmt.Println("Erreur lors de la préparation de la requête SELECT:", err)
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(user.Email, emailTarget)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête SELECT:", err)
		return false
	}
	defer rows.Close()
	if !rows.Next() {
		stmtInsert, err := tx.Prepare("INSERT INTO Follow(AuthorEmail, ReceiverEmail) VALUES(?, ?)")
		if err != nil {
			fmt.Println("Erreur lors de la préparation de la requête INSERT:", err)
			return false
		}
		defer stmtInsert.Close()
		_, err = stmtInsert.Exec(user.Email, emailTarget)
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la requête INSERT:", err)
			return false
		}
		err = tx.Commit()
		if err != nil {
			fmt.Println("Erreur lors du commit de la transaction:", err)
			return false
		}
		return true
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Erreur lors du commit de la transaction:", err)
		return false
	}

	return false
}

func (user *User) Unfollow(emailTarget string) bool {
	tx, err := DB.core.Begin()
	if err != nil {
		fmt.Println("Erreur lors du début de la transaction:", err)
		return false
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	stmt, err := tx.Prepare("DELETE FROM Follow WHERE AuthorEmail = ? AND ReceiverEmail = ?")
	if err != nil {
		fmt.Println("Erreur lors de la préparation de la requête DELETE:", err)
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(user.Email, emailTarget)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête DELETE:", err)
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Erreur lors de la vérification des lignes affectées:", err)
		return false
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("Erreur lors du commit de la transaction:", err)
		return false
	}
	return rowsAffected > 0
}

func (user *User) GetAllFollowers() []User {
	rows, err := DB.core.Query("SELECT AuthorEmail FROM Follow WHERE ReceiverEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var userIn User
		rows.Scan(&userIn.Email)
		userIn, _, _ = DB.GetUser(userIn.Email)
		users = append(users, userIn)
	}
	return users
}

func (user *User) GetAllFollowedAccount() []User {
	rows, err := DB.core.Query("SELECT ReceiverEmail FROM Follow WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var userIn User
		rows.Scan(&userIn.Email)
		userIn, _, _ = DB.GetUser(userIn.Email)
		users = append(users, userIn)
	}
	return users
}

func (User *User) IsLikedPost(id int) bool {
	rows, _ := DB.core.Query("SELECT * FROM Likes WHERE PostId = ? AND AuthorEmail = ?", id, User.Email)
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

func (User *User) IsDislikedPost(id int) bool {
	rows, _ := DB.core.Query("SELECT * FROM Dislikes WHERE PostId = ? AND AuthorEmail = ?", id, User.Email)
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}
