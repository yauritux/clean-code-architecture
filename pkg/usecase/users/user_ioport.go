package users

type UserInputPort interface {
	FetchCurrentUser(id string) (interface{}, error)
	BuildUserUsecaseModel(interface{}) *User
}

type UserOutputPort interface {
}
