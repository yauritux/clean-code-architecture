package inmem

import (
	"github.com/yauritux/cartsvc/pkg/adapter/repository/inmem/model"
	e "github.com/yauritux/cartsvc/pkg/sharedkernel/error"
	uc "github.com/yauritux/cartsvc/pkg/usecase/products"
)

type ProductRepository struct {
	data []*model.Product
}

func NewProductRepository() *ProductRepository {
	productRecords := make([]*model.Product, 0)
	productRecords = append(productRecords, &model.Product{
		ID:    "001",
		Name:  "Shuriken",
		Stock: 1500,
		Price: 250.50,
		Disc:  0.0,
	})
	productRecords = append(productRecords, &model.Product{
		ID:    "002",
		Name:  "Sai",
		Stock: 950,
		Price: 175.25,
		Disc:  0.0,
	})
	return &ProductRepository{data: productRecords}
}

func (r *ProductRepository) FindByProductID(id string) (interface{}, error) {
	if id == "" {
		return nil, e.NewErrNoData("please provide product id")
	}
	for i, p := range r.data {
		if p.ID == id {
			return r.BuildProductUsecaseModel(r.data[i]), nil
		}
	}
	return nil, e.NewErrNoData("no product found for id " + id)
}

func (r *ProductRepository) BuildProductUsecaseModel(prod interface{}) *uc.Product {
	switch prod.(type) {
	case *model.Product:
		u := prod.(*model.Product)
		return &uc.Product{
			ID:    u.ID,
			Name:  u.Name,
			Stock: u.Stock,
			Price: u.Price,
			Disc:  u.Disc,
		}
	default:
		return nil
	}
}
