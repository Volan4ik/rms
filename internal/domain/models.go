package domain

import "time"

type Role struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Description string  `json:"description,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	IsSystem  bool      `json:"is_system"`
}

type Employee struct {
	ID       int64     `json:"id"`
	FullName string    `json:"full_name"`
	Phone    string    `json:"phone"`
	Email    *string   `json:"email,omitempty"`
	RoleID   int64     `json:"role_id"`
	HiredAt  time.Time `json:"hired_at"`
	IsActive bool      `json:"is_active"`
}

type Customer struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Email     *string   `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	VIPLevel  int       `json:"vip_level"`
}

type RestaurantTable struct {
	ID          int64  `json:"id"`
	TableNumber int    `json:"table_number"`
	Seats       int    `json:"seats"`
	IsActive    bool   `json:"is_active"`
	Description string `json:"description,omitempty"`
}

type MenuCategory struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Description string `json:"description,omitempty"`
	SortOrder  int    `json:"sort_order"`
	IsActive   bool   `json:"is_active"`
}

type Product struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Unit        string   `json:"unit"`
	CostPrice   *float64 `json:"cost_price,omitempty"`
	IsAvailable bool     `json:"is_available"`
}

type ProductStock struct {
	ProductID int64     `json:"product_id"`
	Quantity  float64   `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Dish struct {
	ID              int64   `json:"id"`
	CategoryID      int64   `json:"category_id"`
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	CookTimeMinutes int     `json:"cook_time_minutes"`
	IsActive        bool    `json:"is_active"`
	Description     string  `json:"description,omitempty"`
}

type DishIngredient struct {
	ID        int64   `json:"id"`
	DishID    int64   `json:"dish_id"`
	ProductID int64   `json:"product_id"`
	Quantity  float64 `json:"quantity"`
}

type Reservation struct {
	ID           int64     `json:"id"`
	CustomerID   int64     `json:"customer_id"`
	TableID      int64     `json:"table_id"`
	ReservedFrom time.Time `json:"reserved_from"`
	ReservedTo   time.Time `json:"reserved_to"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type Shift struct {
	ID              int64      `json:"id"`
	OpenedBy        int64      `json:"opened_by"`
	ClosedBy        *int64     `json:"closed_by,omitempty"`
	OpenedAt        time.Time  `json:"opened_at"`
	ClosedAt        *time.Time `json:"closed_at,omitempty"`
	Status          string     `json:"status"`
	Note            string     `json:"note,omitempty"`
	ExpectedRevenue *float64   `json:"expected_revenue,omitempty"`
	ActualRevenue   *float64   `json:"actual_revenue,omitempty"`
}

type Order struct {
	ID            int64     `json:"id"`
	TableID       int64     `json:"table_id"`
	CustomerID    *int64    `json:"customer_id,omitempty"`
	WaiterID      int64     `json:"waiter_id"`
	ReservationID *int64    `json:"reservation_id,omitempty"`
	ShiftID       *int64    `json:"shift_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Status        string    `json:"status"`
}

type OrderItem struct {
	ID             int64   `json:"id"`
	OrderID        int64   `json:"order_id"`
	DishID         int64   `json:"dish_id"`
	Quantity       int     `json:"quantity"`
	PriceAtMoment  float64 `json:"price_at_moment"`
	Comment        string  `json:"comment,omitempty"`
}

type Payment struct {
	ID      int64     `json:"id"`
	OrderID int64     `json:"order_id"`
	Amount  float64   `json:"amount"`
	Method  string    `json:"method"`
	PaidAt  time.Time `json:"paid_at"`
	Status  string    `json:"status"`
}

type ImportError struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	Entity       string    `json:"entity"`
	RawData      string    `json:"raw_data"`
	ErrorMessage string    `json:"error_message"`
}

type ShiftRevenue struct {
	ShiftID      int64      `json:"shift_id"`
	OpenedAt     time.Time  `json:"opened_at"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	OrdersCount  int64      `json:"orders_count"`
	TotalRevenue float64    `json:"total_revenue"`
	AvgCheck     *float64   `json:"avg_check,omitempty"`
}

type WaiterPerformance struct {
	WaiterID     int64   `json:"waiter_id"`
	FullName     string  `json:"full_name"`
	OrdersCount  int64   `json:"orders_count"`
	TotalRevenue float64 `json:"total_revenue"`
	AvgCheck     *float64 `json:"avg_check,omitempty"`
}

type DishAvailability struct {
	ID                   int64   `json:"id"`
	Name                 string  `json:"name"`
	Price                float64 `json:"price"`
	IsActive             bool    `json:"is_active"`
	AllProductsAvailable bool    `json:"all_products_available"`
	CanBeOrdered         bool    `json:"can_be_ordered"`
}
