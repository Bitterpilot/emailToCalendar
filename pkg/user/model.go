package user

type UserRepository interface {
	Store(user User)
	FindById(id int) User
}

type User struct {
	Id      int
	Name    string
	Query   string
	Label   string
	IsAdmin bool
}
