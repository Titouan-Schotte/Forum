package dbmanagement

import (
	"fmt"
	"time"
)

func (user *User) GetAllNotifications() []Notification {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT Message,Type,DatePost FROM Notifications WHERE ReceiverEmail = ? ORDER BY DatePost DESC LIMIT ?", user.Email, 5)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var notifications []Notification

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var notification Notification

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&notification.Message, &notification.Type, &notification.Date)
		// Ajout du post à la slice des posts
		notifications = append(notifications, notification)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return notifications
}

func (user *User) GetNotificationById(id int) Notification {
	// Connexion à la base de données

	// Exécution de la requête SQL pour récupérer les posts de la catégorie donnée
	rows, err := DB.core.Query("SELECT Message,Type,DatePost FROM Notifications WHERE ReceiverEmail = ? AND Id = ?", user.Email, id)
	if err != nil {
		return Notification{}
	}
	defer rows.Close()

	// Création d'une slice pour stocker les posts récupérés
	var notifications []Notification

	// Parcours des résultats et création des structures Post
	for rows.Next() {
		var notification Notification

		// Scan des colonnes de la table Post dans les champs correspondants de la structure Post
		rows.Scan(&notification.Message, &notification.Type, &notification.Date)
		// Ajout du post à la slice des posts
		notifications = append(notifications, notification)
	}

	// Vérification des erreurs éventuelles lors du parcours des résultats
	return notifications[0]
}

func (user *User) AddNotification(message string, typeIn string) bool {
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

	stmtInsert, err := tx.Prepare("INSERT INTO Notifications(Message, Type, ReceiverEmail, DatePost) VALUES(?, ?, ?, ?)")
	if err != nil {
		fmt.Println("Erreur lors de la préparation de la requête INSERT:", err)
		return false
	}
	defer stmtInsert.Close()
	// Obtenir la date et l'heure actuelles
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	// Exécuter la requête INSERT
	_, err = stmtInsert.Exec(message, typeIn, user.Email, formattedTime)
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
