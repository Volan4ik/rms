-- Test data generator for Restaurant Management System
-- Uses setseed for reproducibility inside the course context
SELECT setseed(0.42);

-- Base dictionaries
INSERT INTO roles (name, description, is_system)
VALUES 
    ('admin', 'System administrator', TRUE),
    ('manager', 'Shift manager', FALSE),
    ('waiter', 'Service staff', FALSE),
    ('chef', 'Kitchen staff', FALSE)
ON CONFLICT DO NOTHING;

-- Employees
INSERT INTO employees (full_name, phone, email, role_id, hired_at, is_active)
SELECT
    CASE 
        WHEN gs <= 5 THEN format('Manager %s', gs)
        WHEN gs <= 15 THEN format('Chef %s', gs)
        ELSE format('Waiter %s', gs)
    END AS full_name,
    format('+7000%07s', gs) AS phone,
    format('staff%s@rms.local', gs) AS email,
    CASE 
        WHEN gs <= 5 THEN (SELECT id FROM roles WHERE name='manager')
        WHEN gs <= 15 THEN (SELECT id FROM roles WHERE name='chef')
        ELSE (SELECT id FROM roles WHERE name='waiter')
    END AS role_id,
    current_date - (gs % 400),
    TRUE
FROM generate_series(1,60) AS gs
ON CONFLICT DO NOTHING;

-- Restaurant tables
INSERT INTO restaurant_tables (table_number, seats, is_active, description)
SELECT gs, 2 + (gs % 7), TRUE, format('Table near zone %s', gs % 5)
FROM generate_series(1,30) AS gs
ON CONFLICT DO NOTHING;

-- Menu categories
INSERT INTO menu_categories (name, description, sort_order, is_active)
SELECT
    format('Category %s', gs),
    format('Description for category %s', gs),
    gs,
    TRUE
FROM generate_series(1,12) AS gs
ON CONFLICT DO NOTHING;

-- Products and stock (600+)
INSERT INTO products (name, unit, cost_price, is_available)
SELECT
    format('Product %s', gs),
    (ARRAY['kg','l','pcs'])[1 + (gs % 3)],
    ROUND((random() * 200)::numeric, 2),
    TRUE
FROM generate_series(1,650) AS gs
ON CONFLICT DO NOTHING;

INSERT INTO product_stock (product_id, quantity, updated_at)
SELECT p.id, ROUND((random() * 50)::numeric, 3), now() - (p.id % 30) * interval '1 hour'
FROM products p
ON CONFLICT (product_id) DO UPDATE SET quantity = excluded.quantity, updated_at = excluded.updated_at;

-- Customers (800)
INSERT INTO customers (full_name, phone, email, vip_level)
SELECT
    format('Customer %s', gs),
    format('+7111%07s', gs),
    format('customer%s@demo.local', gs),
    (random() * 3)::int
FROM generate_series(1,800) AS gs
ON CONFLICT DO NOTHING;

-- Dishes (550)
INSERT INTO dishes (category_id, name, price, cook_time_minutes, is_active, description)
SELECT
    (SELECT id FROM menu_categories ORDER BY random() LIMIT 1),
    format('Dish %s', gs),
    ROUND((10 + random() * 70)::numeric, 2) * 10,
    5 + (gs % 50),
    TRUE,
    format('Tasty dish number %s', gs)
FROM generate_series(1,550) AS gs
ON CONFLICT DO NOTHING;

-- Dish ingredients (3 per dish)
DO $$
DECLARE
    d RECORD;
    prod_ids BIGINT[];
    p_id BIGINT;
BEGIN
    FOR d IN SELECT id FROM dishes LOOP
        prod_ids := ARRAY(SELECT id FROM products ORDER BY random() LIMIT 3);
        FOREACH p_id IN ARRAY prod_ids LOOP
            INSERT INTO dish_ingredients (dish_id, product_id, quantity)
            VALUES (d.id, p_id, ROUND((0.1 + random() * 3)::numeric, 3))
            ON CONFLICT DO NOTHING;
        END LOOP;
    END LOOP;
END;
$$;

-- Reservations (300)
INSERT INTO reservations (customer_id, table_id, reserved_from, reserved_to, status, created_at)
SELECT
    (SELECT id FROM customers ORDER BY random() LIMIT 1),
    (SELECT id FROM restaurant_tables ORDER BY random() LIMIT 1),
    now() + (gs || ' minutes')::interval,
    now() + (gs || ' minutes')::interval + interval '2 hours',
    (ARRAY['new','confirmed','completed'])[1 + (gs % 3)],
    now() - (gs % 100) * interval '1 hour'
FROM generate_series(1,300) AS gs
ON CONFLICT DO NOTHING;

-- Shifts (60)
INSERT INTO shifts (opened_by, closed_by, opened_at, closed_at, status, note, expected_revenue, actual_revenue)
SELECT
    (SELECT id FROM employees WHERE role_id = (SELECT id FROM roles WHERE name='manager') ORDER BY random() LIMIT 1),
    (SELECT id FROM employees WHERE role_id = (SELECT id FROM roles WHERE name='manager') ORDER BY random() LIMIT 1),
    now() - (gs || ' hours')::interval,
    now() - ((gs - 8) || ' hours')::interval,
    'closed',
    format('Shift %s notes', gs),
    ROUND((500 + random() * 5000)::numeric, 2),
    ROUND((500 + random() * 5000)::numeric, 2)
FROM generate_series(1,60) AS gs
ON CONFLICT DO NOTHING;

-- Orders (5000)
INSERT INTO orders (table_id, customer_id, waiter_id, reservation_id, shift_id, created_at, status)
SELECT
    (SELECT id FROM restaurant_tables ORDER BY random() LIMIT 1),
    (SELECT id FROM customers ORDER BY random() LIMIT 1),
    (SELECT id FROM employees WHERE role_id = (SELECT id FROM roles WHERE name='waiter') ORDER BY random() LIMIT 1),
    NULL,
    (SELECT id FROM shifts ORDER BY random() LIMIT 1),
    now() - (random() * interval '30 days'),
    (ARRAY['new','in_progress','closed','cancelled'])[1 + floor(random()*4)]
FROM generate_series(1,5000);

-- Link some reservations uniquely to orders
WITH ordered_res AS (
    SELECT id, row_number() OVER () AS rn FROM reservations LIMIT 200
)
UPDATE orders o
SET reservation_id = r.id
FROM ordered_res r
WHERE o.id = r.rn;

-- Order items (2-4 per order)
INSERT INTO order_items (order_id, dish_id, quantity, price_at_moment, comment)
SELECT
    o.id,
    d.id,
    1 + (floor(random()*3))::int,
    d.price,
    NULL
FROM orders o
CROSS JOIN LATERAL (
    SELECT id, price FROM dishes ORDER BY random() LIMIT (2 + (o.id % 3))
) d
ON CONFLICT (order_id, dish_id) DO NOTHING;

-- Payments (1:1 with orders)
INSERT INTO payments (order_id, amount, method, status, paid_at)
SELECT
    o.id,
    GREATEST(0, (SELECT SUM(oi.price_at_moment * oi.quantity) FROM order_items oi WHERE oi.order_id = o.id)),
    (ARRAY['cash','card','online'])[1 + floor(random()*3)],
    CASE 
        WHEN o.status = 'cancelled' THEN 'refunded'
        WHEN random() < 0.15 THEN 'pending'
        ELSE 'paid'
    END,
    o.created_at + interval '30 minutes'
FROM orders o
ON CONFLICT (order_id) DO NOTHING;
