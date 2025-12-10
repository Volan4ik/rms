-- Schema for Restaurant Management System
-- Contains tables, constraints, indexes, views, functions, and triggers

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    is_system BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS employees (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    phone TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    role_id BIGINT NOT NULL REFERENCES roles(id),
    hired_at DATE NOT NULL DEFAULT current_date,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    full_name TEXT NOT NULL,
    phone TEXT NOT NULL UNIQUE,
    email TEXT UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    vip_level SMALLINT NOT NULL DEFAULT 0 CHECK (vip_level BETWEEN 0 AND 3)
);

CREATE TABLE IF NOT EXISTS restaurant_tables (
    id BIGSERIAL PRIMARY KEY,
    table_number INT NOT NULL UNIQUE,
    seats INT NOT NULL CHECK (seats > 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS menu_categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    sort_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    unit TEXT NOT NULL,
    cost_price NUMERIC(10,2),
    is_available BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS product_stock (
    product_id BIGINT PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    quantity NUMERIC(12,3) NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dishes (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES menu_categories(id),
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL CHECK (price > 0),
    cook_time_minutes INT NOT NULL CHECK (cook_time_minutes > 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT,
    UNIQUE(category_id, name)
);

CREATE TABLE IF NOT EXISTS dish_ingredients (
    id BIGSERIAL PRIMARY KEY,
    dish_id BIGINT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(id),
    quantity NUMERIC(10,3) NOT NULL CHECK (quantity > 0),
    UNIQUE(dish_id, product_id)
);

CREATE TABLE IF NOT EXISTS reservations (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    table_id BIGINT NOT NULL REFERENCES restaurant_tables(id),
    reserved_from TIMESTAMP NOT NULL,
    reserved_to TIMESTAMP NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('new','confirmed','cancelled','completed')),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CHECK (reserved_to > reserved_from)
);

CREATE TABLE IF NOT EXISTS shifts (
    id BIGSERIAL PRIMARY KEY,
    opened_by BIGINT NOT NULL REFERENCES employees(id),
    closed_by BIGINT REFERENCES employees(id),
    opened_at TIMESTAMP NOT NULL DEFAULT now(),
    closed_at TIMESTAMP,
    status TEXT NOT NULL CHECK (status IN ('opened','closed')),
    note TEXT,
    expected_revenue NUMERIC(12,2),
    actual_revenue NUMERIC(12,2)
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    table_id BIGINT NOT NULL REFERENCES restaurant_tables(id),
    customer_id BIGINT REFERENCES customers(id) ON DELETE SET NULL,
    waiter_id BIGINT NOT NULL REFERENCES employees(id),
    reservation_id BIGINT UNIQUE REFERENCES reservations(id) ON DELETE SET NULL,
    shift_id BIGINT REFERENCES shifts(id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    status TEXT NOT NULL CHECK (status IN ('new','in_progress','closed','cancelled'))
);

CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    dish_id BIGINT NOT NULL REFERENCES dishes(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    price_at_moment NUMERIC(10,2) NOT NULL CHECK (price_at_moment >= 0),
    comment TEXT,
    UNIQUE(order_id, dish_id)
);

CREATE TABLE IF NOT EXISTS payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    amount NUMERIC(10,2) NOT NULL CHECK (amount >= 0),
    method TEXT NOT NULL CHECK (method IN ('cash','card','online')),
    paid_at TIMESTAMP NOT NULL DEFAULT now(),
    status TEXT NOT NULL CHECK (status IN ('pending','paid','refunded'))
);

CREATE TABLE IF NOT EXISTS audit_log (
    id BIGSERIAL PRIMARY KEY,
    table_name TEXT NOT NULL,
    record_id BIGINT NOT NULL,
    operation TEXT NOT NULL CHECK (operation IN ('INSERT','UPDATE','DELETE')),
    changed_at TIMESTAMP NOT NULL DEFAULT now(),
    changed_by TEXT,
    old_data JSONB,
    new_data JSONB
);

CREATE TABLE IF NOT EXISTS import_errors (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    entity TEXT NOT NULL,
    raw_data JSONB,
    error_message TEXT NOT NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_employees_role_id ON employees(role_id);
CREATE INDEX IF NOT EXISTS idx_customers_phone ON customers(phone);
CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(created_at);
CREATE INDEX IF NOT EXISTS idx_tables_active ON restaurant_tables(is_active);
CREATE INDEX IF NOT EXISTS idx_products_available ON products(is_available);
CREATE INDEX IF NOT EXISTS idx_product_stock_product_id ON product_stock(product_id);
CREATE INDEX IF NOT EXISTS idx_dishes_category_id ON dishes(category_id);
CREATE INDEX IF NOT EXISTS idx_dishes_active ON dishes(is_active);
CREATE INDEX IF NOT EXISTS idx_dish_ingredients_dish_id ON dish_ingredients(dish_id);
CREATE INDEX IF NOT EXISTS idx_dish_ingredients_product_id ON dish_ingredients(product_id);
CREATE INDEX IF NOT EXISTS idx_reservations_customer_id ON reservations(customer_id);
CREATE INDEX IF NOT EXISTS idx_reservations_table_id ON reservations(table_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations(status);
CREATE INDEX IF NOT EXISTS idx_reservations_created_at ON reservations(created_at);
CREATE INDEX IF NOT EXISTS idx_shifts_opened_by ON shifts(opened_by);
CREATE INDEX IF NOT EXISTS idx_shifts_closed_by ON shifts(closed_by);
CREATE INDEX IF NOT EXISTS idx_shifts_status ON shifts(status);
CREATE INDEX IF NOT EXISTS idx_orders_table_id ON orders(table_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_waiter_id ON orders(waiter_id);
CREATE INDEX IF NOT EXISTS idx_orders_shift_id ON orders(shift_id);
CREATE INDEX IF NOT EXISTS idx_orders_reservation_id ON orders(reservation_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_dish_id ON order_items(dish_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_paid_at ON payments(paid_at);

-- Functions and triggers

CREATE OR REPLACE FUNCTION fn_audit() RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO audit_log(table_name, record_id, operation, changed_by, old_data, new_data)
    VALUES (TG_TABLE_NAME, COALESCE(NEW.id, OLD.id), TG_OP, CURRENT_USER, to_jsonb(OLD), to_jsonb(NEW));
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_update_product_availability() RETURNS TRIGGER AS $$
DECLARE
    qty NUMERIC(12,3);
BEGIN
    IF TG_OP = 'DELETE' THEN
        qty := 0;
        UPDATE products SET is_available = FALSE WHERE id = OLD.product_id;
        RETURN OLD;
    END IF;

    qty := COALESCE(NEW.quantity, 0);
    UPDATE products SET is_available = (qty > 0) WHERE id = NEW.product_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger bindings for availability
DROP TRIGGER IF EXISTS trg_product_stock_update ON product_stock;
CREATE TRIGGER trg_product_stock_update
AFTER INSERT OR UPDATE OR DELETE ON product_stock
FOR EACH ROW EXECUTE FUNCTION fn_update_product_availability();

-- Audit triggers
DO $$
DECLARE
    tbl TEXT;
BEGIN
    FOREACH tbl IN ARRAY ARRAY['customers','employees','products','dishes','reservations','shifts','orders','order_items','payments'] LOOP
        EXECUTE format('DROP TRIGGER IF EXISTS trg_audit_%s ON %s;', tbl, tbl);
        EXECUTE format('CREATE TRIGGER trg_audit_%s AFTER INSERT OR UPDATE OR DELETE ON %s FOR EACH ROW EXECUTE FUNCTION fn_audit();', tbl, tbl);
    END LOOP;
END;
$$;

-- Views
CREATE OR REPLACE VIEW view_dishes_availability AS
SELECT
    d.id,
    d.name,
    d.price,
    d.is_active,
    CASE WHEN COUNT(di.*) = 0 THEN FALSE ELSE BOOL_AND(p.is_available) END AS all_products_available,
    d.is_active AND (CASE WHEN COUNT(di.*) = 0 THEN FALSE ELSE BOOL_AND(p.is_available) END) AS can_be_ordered
FROM dishes d
LEFT JOIN dish_ingredients di ON di.dish_id = d.id
LEFT JOIN products p ON p.id = di.product_id
GROUP BY d.id, d.name, d.price, d.is_active;

CREATE OR REPLACE VIEW view_shift_revenue AS
SELECT
    s.id AS shift_id,
    s.opened_at,
    s.closed_at,
    COUNT(DISTINCT o.id) AS orders_count,
    COALESCE(SUM(CASE WHEN pay.status = 'paid' THEN pay.amount ELSE 0 END), 0) AS total_revenue,
    CASE WHEN COUNT(DISTINCT o.id) = 0 THEN NULL ELSE ROUND(COALESCE(SUM(CASE WHEN pay.status = 'paid' THEN pay.amount ELSE 0 END), 0) / COUNT(DISTINCT o.id), 2) END AS avg_check
FROM shifts s
LEFT JOIN orders o ON o.shift_id = s.id
LEFT JOIN payments pay ON pay.order_id = o.id
GROUP BY s.id;

CREATE OR REPLACE VIEW view_waiter_performance AS
SELECT
    e.id AS waiter_id,
    e.full_name,
    COUNT(DISTINCT o.id) AS orders_count,
    COALESCE(SUM(CASE WHEN p.status = 'paid' THEN p.amount ELSE 0 END), 0) AS total_revenue,
    CASE WHEN COUNT(DISTINCT o.id) = 0 THEN NULL ELSE ROUND(COALESCE(SUM(CASE WHEN p.status = 'paid' THEN p.amount ELSE 0 END), 0) / COUNT(DISTINCT o.id), 2) END AS avg_check
FROM employees e
LEFT JOIN orders o ON o.waiter_id = e.id
LEFT JOIN payments p ON p.order_id = o.id
WHERE e.is_active = TRUE
GROUP BY e.id, e.full_name;

CREATE OR REPLACE VIEW view_popular_dishes AS
SELECT
    d.id,
    d.name,
    COUNT(oi.id) AS times_ordered,
    SUM(oi.quantity) AS portions_sold,
    COALESCE(SUM(oi.price_at_moment * oi.quantity), 0) AS revenue
FROM dishes d
LEFT JOIN order_items oi ON oi.dish_id = d.id
GROUP BY d.id, d.name
ORDER BY portions_sold DESC NULLS LAST;

-- Scalar functions
CREATE OR REPLACE FUNCTION get_customer_total_spent(p_customer_id BIGINT) RETURNS NUMERIC AS $$
DECLARE
    total NUMERIC;
BEGIN
    SELECT COALESCE(SUM(p.amount), 0) INTO total
    FROM payments p
    JOIN orders o ON o.id = p.order_id
    WHERE o.customer_id = p_customer_id
      AND p.status = 'paid';
    RETURN total;
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION get_shift_revenue(p_shift_id BIGINT) RETURNS NUMERIC AS $$
DECLARE
    total NUMERIC;
BEGIN
    SELECT COALESCE(SUM(p.amount), 0) INTO total
    FROM payments p
    JOIN orders o ON o.id = p.order_id
    WHERE o.shift_id = p_shift_id
      AND p.status = 'paid';
    RETURN total;
END;
$$ LANGUAGE plpgsql STABLE;

-- Table-valued functions
CREATE OR REPLACE FUNCTION get_shift_report(p_from DATE, p_to DATE)
RETURNS TABLE (
    shift_id BIGINT,
    opened_at TIMESTAMP,
    closed_at TIMESTAMP,
    orders_count BIGINT,
    total_revenue NUMERIC,
    avg_check NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        s.id,
        s.opened_at,
        s.closed_at,
        COUNT(DISTINCT o.id) AS orders_count,
        COALESCE(SUM(CASE WHEN pay.status = 'paid' THEN pay.amount ELSE 0 END), 0) AS total_revenue,
        CASE WHEN COUNT(DISTINCT o.id) = 0 THEN NULL ELSE ROUND(COALESCE(SUM(CASE WHEN pay.status = 'paid' THEN pay.amount ELSE 0 END), 0) / COUNT(DISTINCT o.id), 2) END AS avg_check
    FROM shifts s
    LEFT JOIN orders o ON o.shift_id = s.id
    LEFT JOIN payments pay ON pay.order_id = o.id
    WHERE s.opened_at::date BETWEEN p_from AND p_to
    GROUP BY s.id;
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION get_waiter_performance()
RETURNS TABLE (
    waiter_id BIGINT,
    full_name TEXT,
    orders_count BIGINT,
    total_revenue NUMERIC,
    avg_check NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM view_waiter_performance;
END;
$$ LANGUAGE plpgsql STABLE;
