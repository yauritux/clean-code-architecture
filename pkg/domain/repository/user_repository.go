package repository

type UserRepository interface {
	FindByUserID(string) (interface{}, error)
}
