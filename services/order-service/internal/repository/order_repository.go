package repository

import (
	"database/sql"
	"order-service/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *domain.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var orderID int
	err = tx.QueryRow(`INSERT INTO orders(user_id, status) VALUES($1, $2) RETURNING id`, order.UserID, order.Status).Scan(&orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err := tx.Exec(`INSERT INTO order_items(order_id, product_id, quantity) VALUES($1, $2, $3)`, orderID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	order.ID = orderID
	return nil
}

func (r *OrderRepository) GetOrderByID(id int) (*domain.Order, error) {
	order := &domain.Order{}
	err := r.db.QueryRow(`
        SELECT id, user_id, status
        FROM orders
        WHERE id = $1
    `, id).Scan(&order.ID, &order.UserID, &order.Status)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
        SELECT product_id, quantity
        FROM order_items
        WHERE order_id = $1
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.OrderItem
		err := rows.Scan(&item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return order, nil
}

func (r *OrderRepository) UpdateOrderStatus(id int, status string) error {
	_, err := r.db.Exec(`
        UPDATE orders
        SET status = $1
        WHERE id = $2
    `, status, id)
	return err
}
func (r *OrderRepository) GetOrdersByUser(userID int) ([]domain.Order, error) {
	rows, err := r.db.Query(`
        SELECT id, status
        FROM orders
        WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		order.UserID = userID
		err := rows.Scan(&order.ID, &order.Status)
		if err != nil {
			return nil, err
		}

		items, err := r.getOrderItems(order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) getOrderItems(orderID int) ([]domain.OrderItem, error) {
	rows, err := r.db.Query(`
        SELECT product_id, quantity
        FROM order_items
        WHERE order_id = $1
    `, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem
	for rows.Next() {
		var item domain.OrderItem
		err := rows.Scan(&item.ProductID, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
