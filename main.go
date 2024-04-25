package main

import "fmt"

func main() {
	db := LoadDb()
	db.CreateAccount("Titoune", "titouan.scht@gmail.com", "Pass123")

	users, _ := db.GetUsers()
	fmt.Print(users)
}
