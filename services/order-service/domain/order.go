package domain

type Order struct {
    ID        int           
    UserID    int           
    Status    string        
    Items     []OrderItem   
}

type OrderItem struct {
    ID        int 
    OrderID   int 
    ProductID int 
    Quantity  int 
}
