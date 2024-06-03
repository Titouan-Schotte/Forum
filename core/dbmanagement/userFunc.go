package dbmanagement

import (
	"fmt"
	"log"
)

func (db *DBForum) CreateAccount(pseudo string, email string, password string) (User, bool, string) {
	// Préparer la requête d'insertion
	stmt, err := db.core.Prepare("INSERT INTO User(Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Print("Erreur lors de la préparation de la requête d'insertion:", err)
		return User{}, false, "Erreur de la base de données."
	}
	defer stmt.Close()

	// Exécuter la requête d'insertion
	_, err = stmt.Exec(pseudo, email, password, false, false, false, false)
	if err != nil {
		log.Print("Erreur lors de l'exécution de la requête d'insertion:", err)
		return User{}, false, "Erreur de la base de données."
	}

	// Créer un nouvel utilisateur avec les données fournies et l'ID généré
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
	// Préparer la requête de sélection pour récupérer l'utilisateur correspondant à l'email et au mot de passe
	row := db.core.QueryRow("SELECT Pseudo, Email, Password, IsCertified, IsModo, IsAdmin, IsBan FROM User WHERE Email=? AND Password=?", email, password)

	// Créer une structure User pour stocker les données de l'utilisateur
	var user User
	// Scanner les valeurs des colonnes dans la structure User
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

	// Si l'utilisateur est trouvé, retourner la structure User et true
	return user, true, ""
}
func (user *User) EditUser() bool {
	// Vérifier les autorisations de l'utilisateur
	// Préparer la requête de mise à jour du nombre de likes du commentaire
	stmt, err := DB.core.Prepare("UPDATE User SET Pseudo = ?, Email = ?, Password = ?, IsCertified = ?, IsModo = ?, IsAdmin = ?, IsBan = ? WHERE Email = ?")
	if err != nil {
		return false
	}
	defer stmt.Close()

	// Exécuter la requête de mise à jour du nombre de likes du commentaire
	_, err = stmt.Exec(user.Pseudo, user.Email, user.Password, user.IsCertified, user.IsModo, user.IsAdmin, user.IsBan, user.Email)
	if err != nil {
		return false
	}
	// Retourner true pour indiquer que l'édition du post a réussi
	return true
}
func (user *User) DeleteAccount() (bool, string) {
	// Préparer la requête de suppression de l'utilisateur avec l'email et le mot de passe spécifiés
	stmt, err := DB.core.Prepare("DELETE FROM User WHERE Email=? AND Password=?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression:", err)
		return false, "Erreur de la base de données."
	}
	defer stmt.Close()

	// Exécuter la requête de suppression
	_, err = stmt.Exec(user.Email, user.Password)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression:", err)
		return false, "Erreur de la base de données."
	}

	// Si la suppression est réussie, retourner true
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
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT PostId FROM Likes WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var posts []Post

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var post Post

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&post.Id)
		post = DB.GetPostById(user.Email, user.Password, post.Id)
		fmt.Println("idi", post.Title)
		// Ajout du post à la slice des posts
		posts = append(posts, post)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return posts
}

func (user *User) GetAllDislikedPosts() []Post {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT PostId FROM Dislikes WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var posts []Post

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var post Post

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&post.Id)
		post = DB.GetPostById(user.Email, user.Password, post.Id)
		// Ajout du post à la slice des posts
		posts = append(posts, post)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return posts
}

func (user *User) GetAllUserPosts() []Post {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT Id FROM Post WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var posts []Post

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var post Post

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&post.Id)
		fmt.Println(post.Id)
		post = DB.GetPostById(user.Email, user.Password, post.Id)
		// Ajout du post à la slice des posts
		posts = append(posts, post)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return posts
}

func (user *User) Follow(emailTarget string) bool {
	// Démarrer une transaction
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

	// Préparer la requête SELECT
	stmt, err := tx.Prepare("SELECT AuthorEmail, ReceiverEmail FROM Follow WHERE AuthorEmail = ? AND ReceiverEmail = ?")
	if err != nil {
		fmt.Println("Erreur lors de la préparation de la requête SELECT:", err)
		return false
	}
	defer stmt.Close()

	// Exécuter la requête SELECT
	rows, err := stmt.Query(user.Email, emailTarget)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête SELECT:", err)
		return false
	}
	defer rows.Close()

	// Vérifier si l'entrée existe déjà
	if !rows.Next() {
		// Préparer la requête INSERT
		stmtInsert, err := tx.Prepare("INSERT INTO Follow(AuthorEmail, ReceiverEmail) VALUES(?, ?)")
		if err != nil {
			fmt.Println("Erreur lors de la préparation de la requête INSERT:", err)
			return false
		}
		defer stmtInsert.Close()

		// Exécuter la requête INSERT
		_, err = stmtInsert.Exec(user.Email, emailTarget)
		if err != nil {
			fmt.Println("Erreur lors de l'exécution de la requête INSERT:", err)
			return false
		}

		// Commit la transaction
		err = tx.Commit()
		if err != nil {
			fmt.Println("Erreur lors du commit de la transaction:", err)
			return false
		}

		return true
	}

	// Commit la transaction
	err = tx.Commit()
	if err != nil {
		fmt.Println("Erreur lors du commit de la transaction:", err)
		return false
	}

	return false
}

func (user *User) Unfollow(emailTarget string) bool {
	// Démarrer une transaction
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

	// Préparer la requête DELETE
	stmt, err := tx.Prepare("DELETE FROM Follow WHERE AuthorEmail = ? AND ReceiverEmail = ?")
	if err != nil {
		fmt.Println("Erreur lors de la préparation de la requête DELETE:", err)
		return false
	}
	defer stmt.Close()

	// Exécuter la requête DELETE
	result, err := stmt.Exec(user.Email, emailTarget)
	if err != nil {
		fmt.Println("Erreur lors de l'exécution de la requête DELETE:", err)
		return false
	}

	// Vérifier si une ligne a été affectée
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Erreur lors de la vérification des lignes affectées:", err)
		return false
	}

	// Commit la transaction
	err = tx.Commit()
	if err != nil {
		fmt.Println("Erreur lors du commit de la transaction:", err)
		return false
	}

	return rowsAffected > 0
}

func (user *User) GetAllFollowers() []User {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT AuthorEmail FROM Follow WHERE ReceiverEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var users []User

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var userIn User

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&userIn.Email)
		userIn, _, _ = DB.GetUser(userIn.Email)
		// Ajout du post à la slice des posts
		users = append(users, userIn)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return users
}

func (user *User) GetAllFollowedAccount() []User {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT ReceiverEmail FROM Follow WHERE AuthorEmail = ?", user.Email)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var users []User

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var userIn User

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&userIn.Email)
		userIn, _, _ = DB.GetUser(userIn.Email)
		// Ajout du post à la slice des posts
		users = append(users, userIn)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return users
}

func (User *User) IsLikedComment(id int) bool {
	// Connexion à la base de données
	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, _ := DB.core.Query("SELECT * FROM LikesComments WHERE CommentId = ? AND AuthorEmail = ?", id, User.Email)

	defer rows.Close()

	// Parcours des résultats et création des structures Post

	if rows.Next() {
		return true
	}

	return false
}

func (User *User) IsDislikedComment(id int) bool {
	// Connexion à la base de données
	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, _ := DB.core.Query("SELECT * FROM DislikesComments WHERE CommentId = ? AND AuthorEmail = ?", id, User.Email)

	defer rows.Close()

	// Parcours des résultats et création des structures Post

	if rows.Next() {
		return true
	}

	return false
}
func (User *User) IsLikedPost(id int) bool {
	// Connexion à la base de données
	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, _ := DB.core.Query("SELECT * FROM Likes WHERE PostId = ? AND AuthorEmail = ?", id, User.Email)

	defer rows.Close()

	// Parcours des résultats et création des structures Post

	if rows.Next() {
		return true
	}

	return false
}

func (User *User) IsDislikedPost(id int) bool {
	// Connexion à la base de données
	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, _ := DB.core.Query("SELECT * FROM Dislikes WHERE PostId = ? AND AuthorEmail = ?", id, User.Email)

	defer rows.Close()

	// Parcours des résultats et création des structures Post

	if rows.Next() {
		return true
	}

	return false
}
