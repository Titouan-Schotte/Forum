package dbmanagement

import "log"

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

	// Si l'utilisateur est trouvé, retourner la structure User et true
	return user, true, ""
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
