/*
Titouan Schotté

Utilities for users modo management database SQLITE
*/

package dbmanagement

import "log"

func (userModo *User) DeleteAnotherAccountModo(emailTarget string) bool {
	moderator, ok, _ := DB.ConnectToAccount(userModo.Email, userModo.Password)
	if !ok || !(moderator.IsModo || moderator.IsAdmin) {
		// Si le compte n'est pas un modérateur, retourner false
		return false
	}
	targetUser, ok, _ := DB.GetUser(emailTarget)
	if !ok || targetUser.IsModo || targetUser.IsAdmin {
		// Si le compte cible est un modérateur ou n'existe pas, retourner false
		return false
	}
	stmt, err := DB.core.Prepare("DELETE FROM User WHERE Email=?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression:", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(emailTarget)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression:", err)
	}
	return true
}

func (userModo *User) BanAccountModo(emailTarget string) bool {
	moderator, ok, _ := DB.ConnectToAccount(userModo.Email, userModo.Password)
	if !ok || !(moderator.IsModo || moderator.IsAdmin) {
		// Si le compte n'est pas un modérateur, retourner false
		return false
	}
	targetUser, ok, _ := DB.GetUser(emailTarget)
	if !ok || targetUser.IsModo || targetUser.IsAdmin {
		// Si le compte cible est un modérateur ou n'existe pas, retourner false
		return false
	}
	stmt, err := DB.core.Prepare("UPDATE User SET IsBan = ? WHERE Email = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de bannissement:", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(true, emailTarget)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de bannissement:", err)
	}
	return true
}
func (userModo *User) UnbanAccountModo(emailTarget string) bool {
	moderator, ok, _ := DB.ConnectToAccount(userModo.Email, userModo.Password)
	if !ok || !(moderator.IsModo || moderator.IsAdmin) {
		// Si le compte n'est pas un modérateur, retourner false
		return false
	}
	targetUser, ok, _ := DB.GetUser(emailTarget)
	if !ok || targetUser.IsModo || targetUser.IsAdmin {
		// Si le compte cible est un modérateur ou n'existe pas, retourner false
		return false
	}
	stmt, err := DB.core.Prepare("UPDATE User SET IsBan = ? WHERE Email = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de bannissement:", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(false, emailTarget)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de bannissement:", err)
	}
	return true
}
