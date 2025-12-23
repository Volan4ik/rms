package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/example/rms/internal/domain"
)

// Repository contains all DB interactions in one place for simplicity.
type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

var (
	ErrOpenShiftExists = errors.New("shift already opened")
	ErrNoOpenShift     = errors.New("no opened shift")
)

// Common helper
func scanNullableString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func scanNullableFloat(f sql.NullFloat64) *float64 {
	if f.Valid {
		val := f.Float64
		return &val
	}
	return nil
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Customers
func (r *Repository) ListCustomers(ctx context.Context) ([]domain.Customer, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, first_name, last_name, phone, email, created_at, vip_level FROM customers ORDER BY created_at DESC LIMIT 200`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Customer
	for rows.Next() {
		var c domain.Customer
		var email sql.NullString
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Phone, &email, &c.CreatedAt, &c.VIPLevel); err != nil {
			return nil, err
		}
		c.Email = scanNullableString(email)
		res = append(res, c)
	}
	return res, rows.Err()
}

func (r *Repository) CreateCustomer(ctx context.Context, c *domain.Customer) error {
	if err := r.DB.QueryRowContext(ctx, `
		INSERT INTO customers (first_name, last_name, phone, email, vip_level)
		VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`,
		c.FirstName, c.LastName, c.Phone, c.Email, c.VIPLevel).Scan(&c.ID, &c.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateCustomer(ctx context.Context, id int64, c *domain.Customer) error {
	res, err := r.DB.ExecContext(ctx, `
		UPDATE customers SET first_name=$1, last_name=$2, phone=$3, email=$4, vip_level=$5 WHERE id=$6`,
		c.FirstName, c.LastName, c.Phone, c.Email, c.VIPLevel, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteCustomer(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM customers WHERE id=$1`, id)
	return err
}

// Employees
func (r *Repository) ListEmployees(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, first_name, last_name, phone, email, role_id, hired_at, is_active FROM employees ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Employee
	for rows.Next() {
		var e domain.Employee
		var email sql.NullString
		if err := rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.Phone, &email, &e.RoleID, &e.HiredAt, &e.IsActive); err != nil {
			return nil, err
		}
		e.Email = scanNullableString(email)
		res = append(res, e)
	}
	return res, rows.Err()
}

func (r *Repository) CreateEmployee(ctx context.Context, e *domain.Employee) error {
	if err := r.DB.QueryRowContext(ctx, `
		INSERT INTO employees(first_name, last_name, phone, email, role_id, hired_at, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		e.FirstName, e.LastName, e.Phone, e.Email, e.RoleID, e.HiredAt, e.IsActive).Scan(&e.ID); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateEmployee(ctx context.Context, id int64, e *domain.Employee) error {
	res, err := r.DB.ExecContext(ctx, `
		UPDATE employees SET first_name=$1, last_name=$2, phone=$3, email=$4, role_id=$5, is_active=$6 WHERE id=$7`,
		e.FirstName, e.LastName, e.Phone, e.Email, e.RoleID, e.IsActive, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteEmployee(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM employees WHERE id=$1`, id)
	return err
}

// Tables
func (r *Repository) ListTables(ctx context.Context) ([]domain.RestaurantTable, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, table_number, seats, is_active, COALESCE(description,'') FROM restaurant_tables ORDER BY table_number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.RestaurantTable
	for rows.Next() {
		var t domain.RestaurantTable
		if err := rows.Scan(&t.ID, &t.TableNumber, &t.Seats, &t.IsActive, &t.Description); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *Repository) UpsertTable(ctx context.Context, t *domain.RestaurantTable) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO restaurant_tables(table_number, seats, is_active, description)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (table_number) DO UPDATE SET seats=EXCLUDED.seats, is_active=EXCLUDED.is_active, description=EXCLUDED.description
		RETURNING id`, t.TableNumber, t.Seats, t.IsActive, t.Description).Scan(&t.ID)
}

func (r *Repository) DeleteTable(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM restaurant_tables WHERE id=$1`, id)
	return err
}

// Menu categories
func (r *Repository) ListMenuCategories(ctx context.Context) ([]domain.MenuCategory, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, description, sort_order, is_active FROM menu_categories ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.MenuCategory
	for rows.Next() {
		var c domain.MenuCategory
		var desc sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &desc, &c.SortOrder, &c.IsActive); err != nil {
			return nil, err
		}
		if desc.Valid {
			c.Description = desc.String
		}
		res = append(res, c)
	}
	return res, rows.Err()
}

func (r *Repository) UpsertMenuCategory(ctx context.Context, c *domain.MenuCategory) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO menu_categories(name, description, sort_order, is_active)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (name) DO UPDATE SET description=EXCLUDED.description, sort_order=EXCLUDED.sort_order, is_active=EXCLUDED.is_active
		RETURNING id`, c.Name, c.Description, c.SortOrder, c.IsActive).Scan(&c.ID)
}

func (r *Repository) DeleteMenuCategory(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM menu_categories WHERE id=$1`, id)
	return err
}

// Products
func (r *Repository) ListProducts(ctx context.Context, limit int) ([]domain.Product, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, unit, cost_price, is_available FROM products ORDER BY id LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Product
	for rows.Next() {
		var p domain.Product
		var cp sql.NullFloat64
		if err := rows.Scan(&p.ID, &p.Name, &p.Unit, &cp, &p.IsAvailable); err != nil {
			return nil, err
		}
		if cp.Valid {
			val := cp.Float64
			p.CostPrice = &val
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (r *Repository) UpsertProduct(ctx context.Context, p *domain.Product) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO products(name, unit, cost_price, is_available)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (name) DO UPDATE SET unit=EXCLUDED.unit, cost_price=EXCLUDED.cost_price, is_available=EXCLUDED.is_available
		RETURNING id`, p.Name, p.Unit, p.CostPrice, p.IsAvailable).Scan(&p.ID)
}

func (r *Repository) DeleteProduct(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}

// Dishes
func (r *Repository) ListDishes(ctx context.Context, limit int) ([]domain.Dish, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	rows, err := r.DB.QueryContext(ctx, `SELECT id, category_id, name, price, cook_time_minutes, is_active, COALESCE(description,'') FROM dishes ORDER BY id LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Dish
	for rows.Next() {
		var d domain.Dish
		if err := rows.Scan(&d.ID, &d.CategoryID, &d.Name, &d.Price, &d.CookTimeMinutes, &d.IsActive, &d.Description); err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, rows.Err()
}

func (r *Repository) UpsertDish(ctx context.Context, d *domain.Dish) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO dishes(category_id, name, price, cook_time_minutes, is_active, description)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT (category_id, name) DO UPDATE SET price=EXCLUDED.price, cook_time_minutes=EXCLUDED.cook_time_minutes, is_active=EXCLUDED.is_active, description=EXCLUDED.description
		RETURNING id`, d.CategoryID, d.Name, d.Price, d.CookTimeMinutes, d.IsActive, d.Description).Scan(&d.ID)
}

func (r *Repository) DeleteDish(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM dishes WHERE id=$1`, id)
	return err
}

// Reservations
func (r *Repository) ListReservations(ctx context.Context, status string) ([]domain.Reservation, error) {
	query := `SELECT id, customer_id, table_id, reserved_from, reserved_to, status, created_at FROM reservations`
	args := []interface{}{}
	if status != "" {
		query += ` WHERE status = $1`
		args = append(args, status)
	}
	query += ` ORDER BY reserved_from DESC LIMIT 200`

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Reservation
	for rows.Next() {
		var rsv domain.Reservation
		if err := rows.Scan(&rsv.ID, &rsv.CustomerID, &rsv.TableID, &rsv.ReservedFrom, &rsv.ReservedTo, &rsv.Status, &rsv.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, rsv)
	}
	return res, rows.Err()
}

func (r *Repository) CreateReservation(ctx context.Context, rsv *domain.Reservation) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO reservations(customer_id, table_id, reserved_from, reserved_to, status)
		VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`,
		rsv.CustomerID, rsv.TableID, rsv.ReservedFrom, rsv.ReservedTo, rsv.Status).
		Scan(&rsv.ID, &rsv.CreatedAt)
}

func (r *Repository) UpdateReservationStatus(ctx context.Context, id int64, status string) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE reservations SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *Repository) DeleteReservation(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM reservations WHERE id=$1`, id)
	return err
}

// Orders
func (r *Repository) ListOrders(ctx context.Context, status string, limit int) ([]domain.Order, error) {
	if limit <= 0 || limit > 300 {
		limit = 100
	}
	query := `SELECT id, table_id, customer_id, waiter_id, reservation_id, shift_id, created_at, status FROM orders`
	args := []interface{}{}
	if status != "" {
		query += ` WHERE status=$1 ORDER BY created_at DESC LIMIT $2`
		args = append(args, status, limit)
	} else {
		query += ` ORDER BY created_at DESC LIMIT $1`
		args = append(args, limit)
	}

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Order
	for rows.Next() {
		var o domain.Order
		var customer sql.NullInt64
		var reservation sql.NullInt64
		var shift sql.NullInt64
		if err := rows.Scan(&o.ID, &o.TableID, &customer, &o.WaiterID, &reservation, &shift, &o.CreatedAt, &o.Status); err != nil {
			return nil, err
		}
		if customer.Valid {
			val := customer.Int64
			o.CustomerID = &val
		}
		if reservation.Valid {
			val := reservation.Int64
			o.ReservationID = &val
		}
		if shift.Valid {
			val := shift.Int64
			o.ShiftID = &val
		}
		res = append(res, o)
	}
	return res, rows.Err()
}

func (r *Repository) CreateOrder(ctx context.Context, o *domain.Order, items []domain.OrderItem) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = tx.QueryRowContext(ctx, `
		INSERT INTO orders(table_id, customer_id, waiter_id, reservation_id, shift_id, status)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at`,
		o.TableID, o.CustomerID, o.WaiterID, o.ReservationID, o.ShiftID, o.Status).
		Scan(&o.ID, &o.CreatedAt)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO order_items(order_id, dish_id, quantity, price_at_moment, comment)
			VALUES ($1,$2,$3,$4,$5)
			ON CONFLICT (order_id, dish_id) DO UPDATE SET quantity=EXCLUDED.quantity, price_at_moment=EXCLUDED.price_at_moment, comment=EXCLUDED.comment`,
			o.ID, item.DishID, item.Quantity, item.PriceAtMoment, item.Comment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE orders SET status=$1 WHERE id=$2`, status, id)
	return err
}

func (r *Repository) AddOrderItem(ctx context.Context, orderID int64, item domain.OrderItem) error {
	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO order_items(order_id, dish_id, quantity, price_at_moment, comment)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (order_id, dish_id) DO UPDATE SET quantity=EXCLUDED.quantity, price_at_moment=EXCLUDED.price_at_moment, comment=EXCLUDED.comment`,
		orderID, item.DishID, item.Quantity, item.PriceAtMoment, item.Comment)
	return err
}

func (r *Repository) ListOrderItems(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, order_id, dish_id, quantity, price_at_moment, COALESCE(comment,'') FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.OrderItem
	for rows.Next() {
		var oi domain.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.DishID, &oi.Quantity, &oi.PriceAtMoment, &oi.Comment); err != nil {
			return nil, err
		}
		res = append(res, oi)
	}
	return res, rows.Err()
}

func (r *Repository) DeleteOrderItem(ctx context.Context, id int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM order_items WHERE id=$1`, id)
	return err
}

// Payments
func (r *Repository) UpsertPayment(ctx context.Context, p *domain.Payment) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO payments(order_id, amount, method, status, paid_at)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (order_id) DO UPDATE SET amount=EXCLUDED.amount, method=EXCLUDED.method, status=EXCLUDED.status, paid_at=EXCLUDED.paid_at
		RETURNING id`,
		p.OrderID, p.Amount, p.Method, p.Status, p.PaidAt).Scan(&p.ID)
}

func (r *Repository) DeletePayment(ctx context.Context, orderID int64) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM payments WHERE order_id=$1`, orderID)
	return err
}

// Reports
func (r *Repository) GetShiftRevenue(ctx context.Context) ([]domain.ShiftRevenue, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT shift_id, opened_at, closed_at, orders_count, total_revenue, avg_check FROM view_shift_revenue ORDER BY shift_id DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.ShiftRevenue
	for rows.Next() {
		var sr domain.ShiftRevenue
		var closed sql.NullTime
		var avg sql.NullFloat64
		if err := rows.Scan(&sr.ShiftID, &sr.OpenedAt, &closed, &sr.OrdersCount, &sr.TotalRevenue, &avg); err != nil {
			return nil, err
		}
		if closed.Valid {
			val := closed.Time
			sr.ClosedAt = &val
		}
		if avg.Valid {
			val := avg.Float64
			sr.AvgCheck = &val
		}
		res = append(res, sr)
	}
	return res, rows.Err()
}

func (r *Repository) GetWaiterPerformance(ctx context.Context) ([]domain.WaiterPerformance, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT waiter_id, first_name, last_name, orders_count, total_revenue, avg_check FROM view_waiter_performance ORDER BY total_revenue DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.WaiterPerformance
	for rows.Next() {
		var wp domain.WaiterPerformance
		var avg sql.NullFloat64
		if err := rows.Scan(&wp.WaiterID, &wp.FirstName, &wp.LastName, &wp.OrdersCount, &wp.TotalRevenue, &avg); err != nil {
			return nil, err
		}
		if avg.Valid {
			val := avg.Float64
			wp.AvgCheck = &val
		}
		res = append(res, wp)
	}
	return res, rows.Err()
}

func (r *Repository) GetDishesAvailability(ctx context.Context) ([]domain.DishAvailability, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, price, is_active, all_products_available, can_be_ordered FROM view_dishes_availability ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.DishAvailability
	for rows.Next() {
		var da domain.DishAvailability
		if err := rows.Scan(&da.ID, &da.Name, &da.Price, &da.IsActive, &da.AllProductsAvailable, &da.CanBeOrdered); err != nil {
			return nil, err
		}
		res = append(res, da)
	}
	return res, rows.Err()
}

func (r *Repository) GetAvailableDishes(ctx context.Context, categoryID *int64, maxPrice *float64, limit, offset int) ([]domain.DishAvailability, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	query := `
SELECT v.id, v.name, v.price, v.is_active, v.all_products_available, v.can_be_ordered
FROM view_dishes_availability v
JOIN dishes d ON d.id = v.id
WHERE v.can_be_ordered = TRUE`
	args := []interface{}{}
	idx := 1
	if categoryID != nil {
		query += fmt.Sprintf(" AND d.category_id = $%d", idx)
		args = append(args, *categoryID)
		idx++
	}
	if maxPrice != nil {
		query += fmt.Sprintf(" AND v.price <= $%d", idx)
		args = append(args, *maxPrice)
		idx++
	}
	query += fmt.Sprintf(" ORDER BY v.name LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.DishAvailability
	for rows.Next() {
		var da domain.DishAvailability
		if err := rows.Scan(&da.ID, &da.Name, &da.Price, &da.IsActive, &da.AllProductsAvailable, &da.CanBeOrdered); err != nil {
			return nil, err
		}
		res = append(res, da)
	}
	return res, rows.Err()
}

func (r *Repository) ListLowStockProducts(ctx context.Context, threshold float64, limit int) ([]domain.LowStockProduct, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	if threshold <= 0 {
		threshold = 2.0
	}
	rows, err := r.DB.QueryContext(ctx, `
SELECT p.id, p.name, p.unit, COALESCE(ps.quantity, 0) AS quantity, p.is_available
FROM products p
LEFT JOIN product_stock ps ON ps.product_id = p.id
WHERE COALESCE(ps.quantity, 0) < $1
ORDER BY quantity ASC, p.id
LIMIT $2`, threshold, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.LowStockProduct
	for rows.Next() {
		var p domain.LowStockProduct
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Unit, &p.Quantity, &p.IsAvailable); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (r *Repository) GetDishWithIngredients(ctx context.Context, dishID int64) (domain.DishWithIngredients, error) {
	var dish domain.DishWithIngredients
	if err := r.DB.QueryRowContext(ctx, `SELECT id, name, category_id, price FROM dishes WHERE id=$1`, dishID).
		Scan(&dish.DishID, &dish.Name, &dish.CategoryID, &dish.Price); err != nil {
		return dish, err
	}

	rows, err := r.DB.QueryContext(ctx, `
SELECT di.product_id, p.name, di.quantity, p.unit
FROM dish_ingredients di
JOIN products p ON p.id = di.product_id
WHERE di.dish_id = $1`, dishID)
	if err != nil {
		return dish, err
	}
	defer rows.Close()

	for rows.Next() {
		var ing domain.DishIngredientDetail
		if err := rows.Scan(&ing.ProductID, &ing.ProductName, &ing.Quantity, &ing.Unit); err != nil {
			return dish, err
		}
		dish.Ingredients = append(dish.Ingredients, ing)
	}
	return dish, rows.Err()
}

func (r *Repository) ListAvailableTables(ctx context.Context, from, to time.Time, seatsMin int) ([]domain.TableAvailability, error) {
	args := []interface{}{from, to}
	query := `
SELECT t.id, t.table_number, t.seats, t.is_active,
    NOT EXISTS (
        SELECT 1 FROM reservations r
        WHERE r.table_id = t.id
          AND r.status IN ('new','confirmed','completed')
          AND r.reserved_from < $2
          AND r.reserved_to   > $1
    ) AS is_available
FROM restaurant_tables t
WHERE 1=1`
	idx := 3
	if seatsMin > 0 {
		query += fmt.Sprintf(" AND t.seats >= $%d", idx)
		args = append(args, seatsMin)
		idx++
	}
	query += " ORDER BY t.table_number"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.TableAvailability
	for rows.Next() {
		var t domain.TableAvailability
		if err := rows.Scan(&t.TableID, &t.TableNumber, &t.Seats, &t.IsActive, &t.IsAvailable); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *Repository) CheckReservationConflicts(ctx context.Context, tableID int64, from, to time.Time) ([]domain.Reservation, error) {
	rows, err := r.DB.QueryContext(ctx, `
SELECT id, customer_id, table_id, reserved_from, reserved_to, status, created_at
FROM reservations
WHERE table_id = $1
  AND status IN ('new','confirmed','completed')
  AND reserved_from < $3
  AND reserved_to   > $2
ORDER BY reserved_from`, tableID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Reservation
	for rows.Next() {
		var rsv domain.Reservation
		if err := rows.Scan(&rsv.ID, &rsv.CustomerID, &rsv.TableID, &rsv.ReservedFrom, &rsv.ReservedTo, &rsv.Status, &rsv.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, rsv)
	}
	return res, rows.Err()
}

func (r *Repository) OpenShift(ctx context.Context, s *domain.Shift) error {
	if s == nil {
		return errors.New("shift is required")
	}
	err := r.DB.QueryRowContext(ctx, `
INSERT INTO shifts(opened_by, status, note, expected_revenue)
VALUES ($1, 'opened', $2, $3)
RETURNING id, opened_at`, s.OpenedBy, nullIfEmpty(s.Note), s.ExpectedRevenue).
		Scan(&s.ID, &s.OpenedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrOpenShiftExists
		}
		return err
	}
	s.Status = "opened"
	return nil
}

func (r *Repository) CloseShift(ctx context.Context, closedBy int64, note string) (*domain.Shift, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var shift domain.Shift
	if err = tx.QueryRowContext(ctx, `
SELECT id, opened_by, opened_at, note, expected_revenue
FROM shifts
WHERE status='opened'
ORDER BY opened_at DESC
LIMIT 1`).
		Scan(&shift.ID, &shift.OpenedBy, &shift.OpenedAt, &shift.Note, &shift.ExpectedRevenue); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNoOpenShift
		}
		return nil, err
	}

	var revenue sql.NullFloat64
	if err = tx.QueryRowContext(ctx, `
SELECT COALESCE(SUM(p.amount), 0)
FROM payments p
JOIN orders o ON o.id = p.order_id
WHERE o.shift_id = $1 AND p.status = 'paid'`, shift.ID).Scan(&revenue); err != nil {
		return nil, err
	}
	shift.ActualRevenue = scanNullableFloat(revenue)

	if err = tx.QueryRowContext(ctx, `
UPDATE shifts
SET status='closed',
    closed_by=$1,
    closed_at=now(),
    note = COALESCE(NULLIF($2, ''), note),
    actual_revenue=$3
WHERE id=$4
RETURNING closed_at`, closedBy, note, shift.ActualRevenue, shift.ID).
		Scan(&shift.ClosedAt); err != nil {
		return nil, err
	}
	shift.Status = "closed"
	closedByCopy := closedBy
	shift.ClosedBy = &closedByCopy
	return &shift, nil
}

func (r *Repository) GetCurrentShiftOrders(ctx context.Context, status string, waiterID *int64) ([]domain.OrderWithTotal, error) {
	var shiftID int64
	if err := r.DB.QueryRowContext(ctx, `
SELECT id FROM shifts WHERE status='opened' ORDER BY opened_at DESC LIMIT 1`).Scan(&shiftID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoOpenShift
		}
		return nil, err
	}

	query := `
SELECT o.id, o.table_id, o.customer_id, o.waiter_id, o.reservation_id, o.shift_id, o.created_at, o.status,
       COALESCE(SUM(oi.price_at_moment * oi.quantity), 0) AS total_amount
FROM orders o
LEFT JOIN order_items oi ON oi.order_id = o.id
WHERE o.shift_id = $1`
	args := []interface{}{shiftID}
	idx := 2
	if status != "" {
		query += fmt.Sprintf(" AND o.status = $%d", idx)
		args = append(args, status)
		idx++
	}
	if waiterID != nil {
		query += fmt.Sprintf(" AND o.waiter_id = $%d", idx)
		args = append(args, *waiterID)
		idx++
	}
	query += " GROUP BY o.id, o.table_id, o.customer_id, o.waiter_id, o.reservation_id, o.shift_id, o.created_at, o.status ORDER BY o.created_at DESC"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.OrderWithTotal
	for rows.Next() {
		var o domain.OrderWithTotal
		var customer, reservation, shift sql.NullInt64
		if err := rows.Scan(&o.ID, &o.TableID, &customer, &o.WaiterID, &reservation, &shift, &o.CreatedAt, &o.Status, &o.TotalAmount); err != nil {
			return nil, err
		}
		if customer.Valid {
			val := customer.Int64
			o.CustomerID = &val
		}
		if reservation.Valid {
			val := reservation.Int64
			o.ReservationID = &val
		}
		if shift.Valid {
			val := shift.Int64
			o.ShiftID = &val
		}
		res = append(res, o)
	}
	return res, rows.Err()
}

func (r *Repository) AddOrderItemWithPrice(ctx context.Context, orderID, dishID int64, quantity int, comment string) (domain.OrderItem, error) {
	var item domain.OrderItem
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return item, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var status string
	if err = tx.QueryRowContext(ctx, `SELECT status FROM orders WHERE id=$1`, orderID).Scan(&status); err != nil {
		return item, err
	}
	if status == "closed" || status == "cancelled" {
		return item, errors.New("order is not editable")
	}

	if err = tx.QueryRowContext(ctx, `SELECT price FROM dishes WHERE id=$1`, dishID).Scan(&item.PriceAtMoment); err != nil {
		return item, err
	}

	err = tx.QueryRowContext(ctx, `
INSERT INTO order_items(order_id, dish_id, quantity, price_at_moment, comment)
VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (order_id, dish_id) DO UPDATE
SET quantity=EXCLUDED.quantity, price_at_moment=EXCLUDED.price_at_moment, comment=EXCLUDED.comment
RETURNING id`, orderID, dishID, quantity, item.PriceAtMoment, comment).Scan(&item.ID)
	if err != nil {
		return item, err
	}
	item.OrderID = orderID
	item.DishID = dishID
	item.Quantity = quantity
	item.Comment = comment
	return item, nil
}

func (r *Repository) CloseOrder(ctx context.Context, orderID int64, paymentMethod *string) (domain.OrderWithTotal, *domain.Payment, error) {
	var order domain.OrderWithTotal
	var payment *domain.Payment

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return order, nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var customer, reservation, shift sql.NullInt64
	if err = tx.QueryRowContext(ctx, `
SELECT id, table_id, customer_id, waiter_id, reservation_id, shift_id, created_at, status
FROM orders WHERE id=$1 FOR UPDATE`, orderID).
		Scan(&order.ID, &order.TableID, &customer, &order.WaiterID, &reservation, &shift, &order.CreatedAt, &order.Status); err != nil {
		return order, nil, err
	}
	if order.Status == "closed" || order.Status == "cancelled" {
		return order, nil, errors.New("order is already finished")
	}

	if customer.Valid {
		val := customer.Int64
		order.CustomerID = &val
	}
	if reservation.Valid {
		val := reservation.Int64
		order.ReservationID = &val
	}
	if shift.Valid {
		val := shift.Int64
		order.ShiftID = &val
	}

	if err = tx.QueryRowContext(ctx, `
SELECT COALESCE(SUM(price_at_moment * quantity), 0) FROM order_items WHERE order_id=$1`, orderID).Scan(&order.TotalAmount); err != nil {
		return order, nil, err
	}

	if _, err = tx.ExecContext(ctx, `UPDATE orders SET status='closed' WHERE id=$1`, orderID); err != nil {
		return order, nil, err
	}
	order.Status = "closed"

	if paymentMethod != nil && *paymentMethod != "" {
		p := domain.Payment{
			OrderID: orderID,
			Amount:  order.TotalAmount,
			Method:  *paymentMethod,
			Status:  "paid",
			PaidAt:  time.Now(),
		}
		if err = tx.QueryRowContext(ctx, `
INSERT INTO payments(order_id, amount, method, status, paid_at)
VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (order_id) DO UPDATE
SET amount=EXCLUDED.amount, method=EXCLUDED.method, status=EXCLUDED.status, paid_at=EXCLUDED.paid_at
RETURNING id`, p.OrderID, p.Amount, p.Method, p.Status, p.PaidAt).Scan(&p.ID); err != nil {
			return order, nil, err
		}
		payment = &p
	}

	return order, payment, nil
}

func (r *Repository) GetShiftReport(ctx context.Context, from, to time.Time) ([]domain.ShiftRevenue, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT shift_id, opened_at, closed_at, orders_count, total_revenue, avg_check FROM get_shift_report($1,$2) ORDER BY opened_at DESC`, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.ShiftRevenue
	for rows.Next() {
		var sr domain.ShiftRevenue
		var closed sql.NullTime
		var avg sql.NullFloat64
		if err := rows.Scan(&sr.ShiftID, &sr.OpenedAt, &closed, &sr.OrdersCount, &sr.TotalRevenue, &avg); err != nil {
			return nil, err
		}
		if closed.Valid {
			val := closed.Time
			sr.ClosedAt = &val
		}
		if avg.Valid {
			val := avg.Float64
			sr.AvgCheck = &val
		}
		res = append(res, sr)
	}
	return res, rows.Err()
}

func (r *Repository) GetPopularDishes(ctx context.Context, limit int) ([]domain.PopularDish, error) {
	if limit <= 0 || limit > 200 {
		limit = 10
	}
	rows, err := r.DB.QueryContext(ctx, `
SELECT id, name, times_ordered, portions_sold, revenue
FROM view_popular_dishes
ORDER BY portions_sold DESC NULLS LAST
LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.PopularDish
	for rows.Next() {
		var pd domain.PopularDish
		if err := rows.Scan(&pd.DishID, &pd.Name, &pd.TimesOrdered, &pd.PortionsSold, &pd.Revenue); err != nil {
			return nil, err
		}
		res = append(res, pd)
	}
	return res, rows.Err()
}

func (r *Repository) GetCustomerTotalSpent(ctx context.Context, customerID int64) (domain.CustomerTotalSpent, error) {
	var res domain.CustomerTotalSpent
	if err := r.DB.QueryRowContext(ctx, `
SELECT c.id,
       c.first_name,
       c.last_name,
       get_customer_total_spent(c.id)
FROM customers c
WHERE c.id = $1`, customerID).Scan(&res.CustomerID, &res.FirstName, &res.LastName, &res.TotalSpent); err != nil {
		return res, err
	}
	return res, nil
}

func (r *Repository) GetInventoryReport(ctx context.Context) ([]domain.InventoryItem, error) {
	rows, err := r.DB.QueryContext(ctx, `
SELECT p.id, p.name, p.unit, COALESCE(ps.quantity, 0) AS quantity, p.is_available, NOT p.is_available AS in_stop_list
FROM products p
LEFT JOIN product_stock ps ON ps.product_id = p.id
ORDER BY p.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.InventoryItem
	for rows.Next() {
		var it domain.InventoryItem
		if err := rows.Scan(&it.ProductID, &it.Name, &it.Unit, &it.Quantity, &it.IsAvailable, &it.InStopList); err != nil {
			return nil, err
		}
		res = append(res, it)
	}
	return res, rows.Err()
}

func (r *Repository) ListAuditLog(ctx context.Context, tableName string, recordID *int64, from, to *time.Time, limit int) ([]domain.AuditLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	query := `
SELECT id, table_name, record_id, operation, changed_at, changed_by, old_data::text, new_data::text
FROM audit_log
WHERE 1=1`
	args := []interface{}{}
	idx := 1
	if tableName != "" {
		query += fmt.Sprintf(" AND table_name = $%d", idx)
		args = append(args, tableName)
		idx++
	}
	if recordID != nil {
		query += fmt.Sprintf(" AND record_id = $%d", idx)
		args = append(args, *recordID)
		idx++
	}
	if from != nil {
		query += fmt.Sprintf(" AND changed_at >= $%d", idx)
		args = append(args, *from)
		idx++
	}
	if to != nil {
		query += fmt.Sprintf(" AND changed_at <= $%d", idx)
		args = append(args, *to)
		idx++
	}
	query += fmt.Sprintf(" ORDER BY changed_at DESC LIMIT $%d", idx)
	args = append(args, limit)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.AuditLog
	for rows.Next() {
		var a domain.AuditLog
		var changedBy, oldData, newData sql.NullString
		if err := rows.Scan(&a.ID, &a.TableName, &a.RecordID, &a.Operation, &a.ChangedAt, &changedBy, &oldData, &newData); err != nil {
			return nil, err
		}
		if changedBy.Valid {
			a.ChangedBy = &changedBy.String
		}
		if oldData.Valid {
			a.OldData = &oldData.String
		}
		if newData.Valid {
			a.NewData = &newData.String
		}
		res = append(res, a)
	}
	return res, rows.Err()
}

func (r *Repository) ListImportErrors(ctx context.Context, entity string, from, to *time.Time, limit int) ([]domain.ImportError, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	query := `
SELECT id, created_at, entity, raw_data::text, error_message
FROM import_errors
WHERE 1=1`
	args := []interface{}{}
	idx := 1
	if entity != "" {
		query += fmt.Sprintf(" AND entity = $%d", idx)
		args = append(args, entity)
		idx++
	}
	if from != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", idx)
		args = append(args, *from)
		idx++
	}
	if to != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", idx)
		args = append(args, *to)
		idx++
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", idx)
	args = append(args, limit)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.ImportError
	for rows.Next() {
		var ie domain.ImportError
		if err := rows.Scan(&ie.ID, &ie.CreatedAt, &ie.Entity, &ie.RawData, &ie.ErrorMessage); err != nil {
			return nil, err
		}
		res = append(res, ie)
	}
	return res, rows.Err()
}

// Batch import products
func (r *Repository) BatchImportProducts(ctx context.Context, products []domain.Product) (int, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	inserted := 0
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO products(name, unit, cost_price, is_available)
		VALUES ($1,$2,$3,$4)
		ON CONFLICT (name) DO UPDATE SET unit=EXCLUDED.unit, cost_price=EXCLUDED.cost_price, is_available=EXCLUDED.is_available`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, p := range products {
		if p.Name == "" || p.Unit == "" {
			logErr(ctx, tx, "product", p, errors.New("name and unit are required"))
			continue
		}
		if _, err = stmt.ExecContext(ctx, p.Name, p.Unit, p.CostPrice, p.IsAvailable); err != nil {
			logErr(ctx, tx, "product", p, err)
			continue
		}
		inserted++
	}
	return inserted, nil
}

func logErr(ctx context.Context, tx *sql.Tx, entity string, raw interface{}, err error) {
	_, _ = tx.ExecContext(ctx, `
		INSERT INTO import_errors(entity, raw_data, error_message)
		VALUES ($1, to_jsonb($2::json), $3)`,
		entity, rawAsJSON(raw), err.Error())
}

func rawAsJSON(raw interface{}) string {
	b, err := json.Marshal(raw)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func nullableNumber(f *float64) string {
	if f == nil {
		return "null"
	}
	return fmt.Sprintf("%.2f", *f)
}
