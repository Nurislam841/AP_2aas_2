package usecase

import (
    "inventory-service/domain"
    "inventory-service/internal/repository"
)

type ProductUsecase struct {
    repo *repository.ProductRepository
}

func NewProductUsecase(repo *repository.ProductRepository) *ProductUsecase {
    return &ProductUsecase{repo: repo}
}

func (uc *ProductUsecase) GetProduct(id int32) (*domain.Product, error) {
    return uc.repo.GetProductByID(id)
}

func (uc *ProductUsecase) CreateProduct(product domain.Product) (domain.Product, error) {
	id, err := uc.repo.CreateProduct(&product)
	if err != nil {
		return domain.Product{}, err
	}
	product.ID = int(id)
	return product, nil
}


func (uc *ProductUsecase) UpdateProduct(product *domain.Product) error {
	return uc.repo.UpdateProduct(product)
}

func (uc *ProductUsecase) DeleteProduct(id int32) error {
	return uc.repo.DeleteProduct(id)
}
func (uc *ProductUsecase) GetAllProducts(name string, category, limit, offset int) ([]domain.Product, error) {
    return uc.repo.GetAllProducts(name, category, limit, offset)
}
