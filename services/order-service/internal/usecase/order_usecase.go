package usecase

import (
	"order-service/domain"
	"order-service/internal/repository"
	"proto/orderpb"
)

type OrderUsecase struct {
	repo *repository.OrderRepository
}

func NewOrderUsecase(repo *repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}
func (u *OrderUsecase) CreateOrder(order *domain.Order) error {
	order.Status = "pending"
	return u.repo.CreateOrder(order)
}
func (uc *OrderUsecase) GetOrderByID(id int) (*domain.Order, error) {
	return uc.repo.GetOrderByID(id)
}
func (uc *OrderUsecase) UpdateOrderStatus(id int, status string) error {
	return uc.repo.UpdateOrderStatus(id, status)
}
func (uc *OrderUsecase) GetOrdersByUser(userID int) ([]domain.Order, error) {
	return uc.repo.GetOrdersByUser(userID)
}
func (u *OrderUsecase) CreateOrderFromRequest(req *orderpb.CreateOrderRequest) (*domain.Order, error) {
	order := &domain.Order{
		UserID: int(req.UserId),
		Status: "pending",
	}
	for _, i := range req.Items {
		order.Items = append(order.Items, domain.OrderItem{
			ProductID: int(i.ProductId),
			Quantity:  int(i.Quantity),
		})
	}
	err := u.repo.CreateOrder(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (u *OrderUsecase) GetOrderByIDPB(id int) (*orderpb.Order, error) {
	order, err := u.repo.GetOrderByID(id)
	if err != nil {
		return nil, err
	}
	items := []*orderpb.OrderItem{}
	for _, i := range order.Items {
		items = append(items, &orderpb.OrderItem{
			ProductId: int32(i.ProductID),
			Quantity:  int32(i.Quantity),
		})
	}
	return &orderpb.Order{
		Id:     int32(order.ID),
		UserId: int32(order.UserID),
		Status: order.Status,
		Items:  items,
	}, nil
}
func (u *OrderUsecase) GetOrdersByUserPB(userID int) ([]*orderpb.Order, error) {
	orders, err := u.repo.GetOrdersByUser(userID)
	if err != nil {
		return nil, err
	}
	var pbOrders []*orderpb.Order
	for _, o := range orders {
		items := []*orderpb.OrderItem{}
		for _, i := range o.Items {
			items = append(items, &orderpb.OrderItem{
				ProductId: int32(i.ProductID),
				Quantity:  int32(i.Quantity),
			})
		}
		pbOrders = append(pbOrders, &orderpb.Order{
			Id:     int32(o.ID),
			UserId: int32(o.UserID),
			Status: o.Status,
			Items:  items,
		})
	}
	return pbOrders, nil
}
