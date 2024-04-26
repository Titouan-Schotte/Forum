Créer un compte

```go

	//Créer un compte
	user, success, errorMsg := dbmanagement.DB.CreateAccount("Titou", "titou@gmail.com", "Pass123")
	if success {
		fmt.Println("Create :", user.Pseudo)
	} else {
		fmt.Println("Error create :", errorMsg)
	}
```


Connecter un compte


```go
	//Connecter un compte + récup les infos du compte créé
	user, success, errorMsg = dbmanagement.DB.ConnectToAccount("mich@gmail.com", "Pass123")
	if success {
		fmt.Println("Connect :", user.Pseudo)
	} else {
		fmt.Println("Error connect :", errorMsg)
	}
```


Delete un compte

```go
	success, errorMsg = user.DeleteAccount()
	if success {
		fmt.Println("Delete !")
	} else {
		fmt.Println("Error delete :", errorMsg)
	}
```

Afficher tous les comptes

```go
	users, error := dbmanagement.DB.GetUsers()
	for i := 0; i < len(users); i++ {
		user := users[i]
		fmt.Println(user.Email)
	}
```


> MODO

Delete un compte modo
```go
	user.DeleteAnotherAccountModo("mich@gmail.com")

```

Ban un compte modo
```go
	user.BanAnotherAccountModo("mich@gmail.com")

```
