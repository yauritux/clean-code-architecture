package usecase

import (
	"fmt"

	"github.com/yauritux/cartsvc/pkg/domain/repository"
)

type ProductUsecase struct {
	repo repository.ProductRepository
}

type Product struct {
	ID    string
	Name  string
	Stock int
	Price float64
	Disc  float64
}

func NewProductUsecase(r repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{r}
}

func (prod *ProductUsecase) FindByProductID(id string) (interface{}, error) {
	p, err := prod.repo.FindByProductID(id)
	if err != nil {
		return nil, err
	}

	productFound, ok := p.(*Product)
	if !ok {
		return nil, fmt.Errorf("cannot find product with ID %s, got an invalid product type returned from the repository", id)
	}

	return productFound, nil
}
