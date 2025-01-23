package repository

type AuthByEmail struct {
	UserId       string
	HashPassword string
	Phone        string
}

type AuthByPhone struct {
	UserId       string
	HashPassword string
	Email        string
}
