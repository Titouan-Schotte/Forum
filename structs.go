package Forum

type Categorie struct {
	Id    int
	Nom   string
	Posts []Post
}

type Comment struct {
	Id      int
	Content string
	Author  User
	Like    int
	Dislike int
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
	Categories  []Categorie
	Author      User
	Comments    []Comment
}

type User struct {
	Id          int
	Pseudo      string
	Email       string
	Password    string
	Posts       []Post
	IsCertified bool
	IsModo      bool
	IsAdmin     bool
	IsBan       bool
}
