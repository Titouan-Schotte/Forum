package dbmanagement

type Categorie struct {
	Id             int
	Nom            string
	Posts          []Post
	UserManipuling User
}

type Comment struct {
	Id         int
	Content    string
	Author     User
	PostOrigin Post
	Like       int
	Dislike    int
}

type Post struct {
	Id          int
	Title       string
	Description string
	Photos      []string
	Danger      int
	Beauty      int
	Like        int
	Dislike     int
	Categorie   Categorie
	Author      User
	Comments    []Comment
	AuthorEmail string
}

type User struct {
	Pseudo      string
	Email       string
	Password    string
	IsCertified bool
	IsModo      bool
	IsAdmin     bool
	IsBan       bool
}
