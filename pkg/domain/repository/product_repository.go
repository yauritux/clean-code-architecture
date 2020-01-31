package repository

type ProductRepository interface {
	FindByProductID(string) (interface{}, error)
}
