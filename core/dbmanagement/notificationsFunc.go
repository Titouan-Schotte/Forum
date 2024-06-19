/*
Titouan Schotté

Utilities for notifications management database SQLITE
*/

package dbmanagement

import (
	"fmt"
	"time"
)

func (user *User) GetAllNotifications() []Notification {

	rows, err := DB.core.Query("SELECT Message,Type,DatePost FROM Notifications WHERE ReceiverEmail = ? ORDER BY DatePost DESC LIMIT ?", user.Email, 5)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notification Notification

		rows.Scan(&notification.Message, &notification.Type, &notification.Date)
		notifications = append(notifications, notification)
	}

	return notifications
}

func (user *User) GetNotificationById(id int) Notification {
	rows, err := DB.core.Query("SELECT Message,Type,DatePost FROM Notifications WHERE ReceiverEmail = ? AND Id = ?", user.Email, id)
	if err != nil {
		return Notification{}
	}
	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var notification Notification

		rows.Scan(&notification.Message, &notification.Type, &notification.Date)
		notifications = append(notifications, notification)
	}

	return notifications[0]
}

func (user *User) AddNotification(message string, typeIn string) bool {
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
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")

	_, err = stmtInsert.Exec(message, typeIn, user.Email, formattedTime)
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
