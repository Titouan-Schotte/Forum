package dbmanagement

import "log"

func (userModo *User) DeleteAnotherAccountModo(emailTarget string) bool {
	// Vérifier si le compte qui initie l'action est un modérateur
	moderator, ok, _ := DB.ConnectToAccount(userModo.Email, userModo.Password)
	if !ok || !(moderator.IsModo || moderator.IsAdmin) {
		// Si le compte n'est pas un modérateur, retourner false
		return false
	}

	// Vérifier si le compte cible n'est pas un modérateur
	targetUser, ok, _ := DB.GetUser(emailTarget)
	if !ok || targetUser.IsModo || targetUser.IsAdmin {
		// Si le compte cible est un modérateur ou n'existe pas, retourner false
		return false
	}

	// Préparer la requête de suppression du compte cible
	stmt, err := DB.core.Prepare("DELETE FROM User WHERE Email=?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de suppression:", err)
	}
	defer stmt.Close()

	// Exécuter la requête de suppression
	_, err = stmt.Exec(emailTarget)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de suppression:", err)
	}

	// Si la suppression est réussie, retourner true
	return true
}

func (userModo *User) BanAccountModo(emailTarget string) bool {
	// Vérifier si le compte qui initie l'action est un modérateur
	moderator, ok, _ := DB.ConnectToAccount(userModo.Email, userModo.Password)
	if !ok || !(moderator.IsModo || moderator.IsAdmin) {
		// Si le compte n'est pas un modérateur, retourner false
		return false
	}

	// Vérifier si le compte cible n'est pas un modérateur
	targetUser, ok, _ := DB.GetUser(emailTarget)
	if !ok || targetUser.IsModo || targetUser.IsAdmin {
		// Si le compte cible est un modérateur ou n'existe pas, retourner false
		return false
	}

	// Préparer la requête de mise à jour du statut de bannissement du compte cible
	stmt, err := DB.core.Prepare("UPDATE User SET IsBan = ? WHERE Email = ?")
	if err != nil {
		log.Fatal("Erreur lors de la préparation de la requête de bannissement:", err)
	}
	defer stmt.Close()

	// Exécuter la requête de bannissement
	_, err = stmt.Exec(true, emailTarget)
	if err != nil {
		log.Fatal("Erreur lors de l'exécution de la requête de bannissement:", err)
	}

	// Si le bannissement est réussi, retourner true
	return true
}
