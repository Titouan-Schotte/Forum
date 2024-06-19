/*
Titouan Schotté

Utilities for users admin management database SQLITE
*/

package dbmanagement

import "log"

func (userAdmin *User) SetAnotherAccountToModo(emailTarget string, state bool) bool {
	admin, ok, _ := DB.ConnectToAccount(userAdmin.Email, userAdmin.Password)
	if !ok || !admin.IsAdmin {
		// Si le compte n'est pas un administrateur, retourner false
		return false
	}
	targetUser, ok, _ := DB.GetUser(emailTarget)
	if !ok || targetUser.IsModo {
		// Si le compte cible est un modérateur ou n'existe pas, retourner false
		return false
	}
	stmt, err := DB.core.Prepare("UPDATE User SET IsModo = ? WHERE Email = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de définition du modérateur:", err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(state, emailTarget)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de définition du modérateur:", err)
	}
	return true
}
func (userAdmin *User) RankModo(emailTarget string) bool {
	return userAdmin.SetAnotherAccountToModo(emailTarget, true)
}
func (userAdmin *User) UnrankModo(emailTarget string) bool {
	return userAdmin.SetAnotherAccountToModo(emailTarget, false)
}
