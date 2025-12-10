package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/example/rms/internal/domain"
)

// Repository contains all DB interactions in one place for simplicity.
type Repository struct {
	DB *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

// Common helper
func scanNullableString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

// Customers
func (r *Repository) ListCustomers(ctx context.Context) ([]domain.Customer, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, full_name, phone, email, created_at, vip_level FROM customers ORDER BY created_at DESC LIMIT 200`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.Customer
	for rows.Next() {
		var c domain.Customer
		var email sql.NullString
		if err := rows.Scan(&c.ID, &c.FullName, &c.Phone, &email, &c.CreatedAt, &c.VIPLevel); err != nil {
			return nil, err
		}
		c.Email = scanNullableString(email)
		res = append(res, c)
	}
	return res, rows.Err()
}

func (r *Repository) CreateCustomer(ctx context.Context, c *domain.Customer) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO customers (full_name, phone, email, vip_level)
		VALUES ($1,$2,$3,$4) RETURNING id, created_at`,
		c.FullName, c.Phone, c.Email, c.VIPLevel).Scan(&c.ID, &c.CreatedAt)
}

func (r *Repository) UpdateCustomer(ctx context.Context, id int64, c *domain.Customer) error {
	res, err := r.DB.ExecContext(ctx, `
		UPDATE customers SET full_name=$1, phone=$2, email=$3, vip_level=$4 WHERE id=$5`,
		c.FullName, c.Phone, c.Email, c.VIPLevel, id)
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
	rows, err := r.DB.QueryContext(ctx, `SELECT id, full_name, phone, email, role_id, hired_at, is_active FROM employees ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.Employee
	for rows.Next() {
		var e domain.Employee
		var email sql.NullString
		if err := rows.Scan(&e.ID, &e.FullName, &e.Phone, &email, &e.RoleID, &e.HiredAt, &e.IsActive); err != nil {
			return nil, err
		}
		e.Email = scanNullableString(email)
		res = append(res, e)
	}
	return res, rows.Err()
}

func (r *Repository) CreateEmployee(ctx context.Context, e *domain.Employee) error {
	return r.DB.QueryRowContext(ctx, `
		INSERT INTO employees(full_name, phone, email, role_id, hired_at, is_active)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		e.FullName, e.Phone, e.Email, e.RoleID, e.HiredAt, e.IsActive).Scan(&e.ID)
}

func (r *Repository) UpdateEmployee(ctx context.Context, id int64, e *domain.Employee) error {
	res, err := r.DB.ExecContext(ctx, `
		UPDATE employees SET full_name=$1, phone=$2, email=$3, role_id=$4, is_active=$5 WHERE id=$6`,
		e.FullName, e.Phone, e.Email, e.RoleID, e.IsActive, id)
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
	rows, err := r.DB.QueryContext(ctx, `SELECT waiter_id, full_name, orders_count, total_revenue, avg_check FROM view_waiter_performance ORDER BY total_revenue DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []domain.WaiterPerformance
	for rows.Next() {
		var wp domain.WaiterPerformance
		var avg sql.NullFloat64
		if err := rows.Scan(&wp.WaiterID, &wp.FullName, &wp.OrdersCount, &wp.TotalRevenue, &avg); err != nil {
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
