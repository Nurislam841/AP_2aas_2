package repository

import (
	"database/sql"
	"fmt"
	"inventory-service/domain"
)

type ProductRepository struct {
    db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductByID(id int32) (*domain.Product, error) {
    product := &domain.Product{}
    err := r.db.QueryRow("SELECT id, name, price, stock, category_id FROM public.products WHERE id = $1", id).
        Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryID)
    if err != nil {
        return nil, err
    }
    return product, nil
}

func (r *ProductRepository) CreateProduct(product *domain.Product) (int32, error) {
    var id int32
	err := r.db.QueryRow(`
        INSERT INTO products (name, price, stock, category_id) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id
    `, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&id)
    fmt.Println(id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProductRepository) UpdateProduct(product *domain.Product) error {
    _, err := r.db.Exec("UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5",
        product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
    return err
}

func (r *ProductRepository) DeleteProduct(id int32) error {
    _, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
    return err
}

func (r *ProductRepository) GetAllProducts(name string, category, limit, offset int) ([]domain.Product, error) {
    query := `SELECT id, name, price, stock, category_id
              FROM products
			  WHERE ($1 = 0 OR category_id = $1)
			  LIMIT $2 OFFSET $3;
              `
    rows, err := r.db.Query(query, category, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []domain.Product
    for rows.Next() {
        var p domain.Product
        err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)
        if err != nil {
            return nil, err
        }
        products = append(products, p)
    }

    return products, nil
}