-- Реалистичный генератор тестовых данных для Restaurant Management System
-- Фиксируем seed, чтобы данные были воспроизводимыми
SELECT setseed(0.42);

-----------
-- ROLES --
-----------
INSERT INTO roles (name, description, is_system)
VALUES 
    ('admin',      'Администратор системы', TRUE),
    ('manager',    'Управляющий сменой',    FALSE),
    ('waiter',     'Официант',              FALSE),
    ('chef',       'Повар',                 FALSE),
    ('bartender',  'Бармен',                FALSE)
ON CONFLICT (name) DO NOTHING;

---------------
-- EMPLOYEES --
---------------
WITH first_names AS (
    SELECT unnest(ARRAY[
        'Иван','Мария','Алексей','Анна','Дмитрий','Екатерина',
        'Сергей','Ольга','Павел','Наталья','Михаил','Елена'
    ]) AS first_name
),
last_names AS (
    SELECT unnest(ARRAY[
        'Иванов','Петров','Сидоров','Кузнецов','Смирнов',
        'Васильев','Попов','Новиков','Фёдоров','Морозов',
        'Волкова','Соколова'
    ]) AS last_name
)
INSERT INTO employees (full_name, phone, email, role_id, hired_at, is_active)
SELECT
    fn.first_name || ' ' || ln.last_name AS full_name,
    format('+79%08s', gs) AS phone,
    format('staff%02s@rms.local', gs) AS email,
    CASE
        WHEN gs <= 2  THEN (SELECT id FROM roles WHERE name='admin')
        WHEN gs <= 6  THEN (SELECT id FROM roles WHERE name='manager')
        WHEN gs <= 18 THEN (SELECT id FROM roles WHERE name='chef')
        WHEN gs <= 24 THEN (SELECT id FROM roles WHERE name='bartender')
        ELSE             (SELECT id FROM roles WHERE name='waiter')
    END AS role_id,
    current_date - ((gs % 730)) * interval '1 day' AS hired_at,
    TRUE
FROM generate_series(1,40) AS gs
CROSS JOIN LATERAL (
    SELECT first_name FROM first_names ORDER BY random() LIMIT 1
) fn
CROSS JOIN LATERAL (
    SELECT last_name FROM last_names ORDER BY random() LIMIT 1
) ln
ON CONFLICT (phone) DO NOTHING;

-----------------------
-- RESTAURANT TABLES --
-----------------------
INSERT INTO restaurant_tables (table_number, seats, is_active, description)
SELECT
    gs AS table_number,
    CASE
        WHEN gs <= 10 THEN 2
        WHEN gs <= 20 THEN 4
        ELSE 6
    END AS seats,
    TRUE,
    CASE
        WHEN gs <= 10 THEN 'Столик у окна'
        WHEN gs <= 20 THEN 'Столик в центре зала'
        ELSE 'Столик возле бара'
    END AS description
FROM generate_series(1,30) AS gs
ON CONFLICT (table_number) DO NOTHING;

---------------------
-- MENU CATEGORIES --
---------------------
INSERT INTO menu_categories (name, description, sort_order, is_active)
VALUES
    ('Салаты',               'Холодные и тёплые салаты',                 1, TRUE),
    ('Супы',                 'Горячие и холодные супы',                  2, TRUE),
    ('Горячие блюда',        'Основные горячие блюда',                   3, TRUE),
    ('Гарниры',              'Гарниры к основным блюдам',                4, TRUE),
    ('Десерты',              'Сладкие блюда',                             5, TRUE),
    ('Безалкогольные напитки','Соки, лимонады, вода',                    6, TRUE),
    ('Алкогольные напитки',  'Вино, пиво, крепкие напитки',              7, TRUE),
    ('Завтраки',             'Завтраки до полудня',                       8, TRUE)
ON CONFLICT (name) DO NOTHING;

--------------
-- PRODUCTS --
--------------
WITH base_products AS (
    SELECT * FROM (VALUES
        ('Куриное филе','kg', 320.00),
        ('Говядина мраморная','kg', 850.00),
        ('Свинина шейка','kg', 420.00),
        ('Лосось охлаждённый','kg', 1100.00),
        ('Картофель','kg', 40.00),
        ('Морковь','kg', 30.00),
        ('Лук репчатый','kg', 25.00),
        ('Помидоры','kg', 150.00),
        ('Огурцы','kg', 140.00),
        ('Салат айсберг','kg', 260.00),
        ('Сливочное масло','kg', 600.00),
        ('Растительное масло','l', 180.00),
        ('Сливки 33%','l', 260.00),
        ('Молоко','l', 70.00),
        ('Яйцо куриное','pcs', 10.00),
        ('Сыр твёрдый','kg', 650.00),
        ('Сыр фета','kg', 720.00),
        ('Сыр моцарелла','kg', 700.00),
        ('Спагетти','kg', 120.00),
        ('Рис','kg', 90.00),
        ('Булочка для бургера','pcs', 15.00),
        ('Котлета говяжья','pcs', 80.00),
        ('Куриное крыло','pcs', 35.00),
        ('Креветки очищенные','kg', 950.00),
        ('Лимон','pcs', 20.00),
        ('Сахар','kg', 60.00),
        ('Мука','kg', 55.00),
        ('Мёд','kg', 500.00),
        ('Кофе зерновой','kg', 900.00),
        ('Чай чёрный','kg', 400.00),
        ('Минеральная вода 0.5л','pcs', 35.00),
        ('Кола 0.5л','pcs', 50.00),
        ('Апельсиновый сок 1л','pcs', 110.00),
        ('Красное вино','l', 750.00),
        ('Белое вино','l', 730.00),
        ('Пиво разливное','l', 140.00),
        ('Водка','l', 600.00),
        ('Виски','l', 1500.00)
    ) AS t(name, unit, cost_price)
)
INSERT INTO products (name, unit, cost_price, is_available)
SELECT
    CASE
        WHEN gs <= (SELECT count(*) FROM base_products)
            THEN (SELECT name FROM base_products ORDER BY name LIMIT 1 OFFSET gs-1)
        ELSE (SELECT name FROM base_products ORDER BY random() LIMIT 1) || ' партия ' || gs
    END AS name,
    (SELECT unit FROM base_products ORDER BY random() LIMIT 1) AS unit,
    (
        (SELECT cost_price FROM base_products ORDER BY random() LIMIT 1)
        * (0.8 + random()*0.4)
    )::numeric(10,2) AS cost_price,
    TRUE
FROM generate_series(1,200) AS gs
ON CONFLICT (name) DO NOTHING;

-------------------
-- PRODUCT STOCK --
-------------------
INSERT INTO product_stock (product_id, quantity, updated_at)
SELECT
    p.id,
    CASE
        WHEN random() < 0.2 THEN 0       -- часть продуктов закончилась
        ELSE ROUND((1 + random() * 50)::numeric, 3)
    END AS quantity,
    now() - (p.id % 72) * interval '1 hour'
FROM products p
ON CONFLICT (product_id) DO UPDATE
SET quantity   = EXCLUDED.quantity,
    updated_at = EXCLUDED.updated_at;

-- Триггер автоматически обновит products.is_available в зависимости от quantity

---------------
-- CUSTOMERS --
---------------
WITH first_names AS (
    SELECT unnest(ARRAY[
        'Александр','Максим','Илья','Кирилл','Тимофей',
        'София','Алиса','Полина','Варвара','Анастасия'
    ]) AS first_name
),
last_names AS (
    SELECT unnest(ARRAY[
        'Иванов','Петров','Сидоров','Кузнецов','Смирнов',
        'Орлова','Сергеева','Попова','Алексеева','Михайлова'
    ]) AS last_name
)
INSERT INTO customers (full_name, phone, email, vip_level, created_at)
SELECT
    fn.first_name || ' ' || ln.last_name AS full_name,
    format('+79%08s', gs + 1000) AS phone,
    format('customer%03s@example.local', gs) AS email,
    CASE
        WHEN random() < 0.7  THEN 0
        WHEN random() < 0.9  THEN 1
        WHEN random() < 0.98 THEN 2
        ELSE                      3
    END AS vip_level,
    now() - (random() * interval '180 days')
FROM generate_series(1,800) AS gs
CROSS JOIN LATERAL (
    SELECT first_name FROM first_names ORDER BY random() LIMIT 1
) fn
CROSS JOIN LATERAL (
    SELECT last_name FROM last_names ORDER BY random() LIMIT 1
) ln
ON CONFLICT (phone) DO NOTHING;

------------
-- DISHES --
------------
WITH cat AS (
    SELECT id, name FROM menu_categories
),
base_dishes AS (
    SELECT 'Салаты'::text AS category_name, unnest(ARRAY[
        'Салат Цезарь с курицей',
        'Греческий салат',
        'Овощной салат с маслом',
        'Тёплый салат с говядиной'
    ]) AS dish_name, 420::int AS base_price, 15 AS base_cook
    UNION ALL
    SELECT 'Супы', unnest(ARRAY[
        'Борщ классический',
        'Суп-пюре из тыквы',
        'Куриный суп с лапшой',
        'Солянка мясная'
    ]), 350, 20
    UNION ALL
    SELECT 'Горячие блюда', unnest(ARRAY[
        'Стейк из говядины',
        'Куриное филе гриль',
        'Паста Карбонара',
        'Ризотто с грибами',
        'Бургер классический'
    ]), 650, 25
    UNION ALL
    SELECT 'Гарниры', unnest(ARRAY[
        'Картофель фри',
        'Картофельное пюре',
        'Овощи гриль',
        'Рис с овощами'
    ]), 220, 10
    UNION ALL
    SELECT 'Десерты', unnest(ARRAY[
        'Чизкейк Нью-Йорк',
        'Шоколадный фондан',
        'Тирамису',
        'Мороженое с фруктами'
    ]), 380, 10
    UNION ALL
    SELECT 'Безалкогольные напитки', unnest(ARRAY[
        'Домашний лимонад',
        'Морс ягодный',
        'Эспрессо',
        'Капучино',
        'Американо'
    ]), 180, 5
    UNION ALL
    SELECT 'Алкогольные напитки', unnest(ARRAY[
        'Бокал красного вина',
        'Бокал белого вина',
        'Пиво светлое разливное',
        'Пиво тёмное разливное'
    ]), 320, 5
    UNION ALL
    SELECT 'Завтраки', unnest(ARRAY[
        'Омлет с ветчиной и сыром',
        'Сырники со сметаной',
        'Каша овсяная на молоке'
    ]), 300, 10
)
INSERT INTO dishes (category_id, name, price, cook_time_minutes, is_active, description)
SELECT
    c.id AS category_id,
    bd.dish_name || CASE WHEN n > 1 THEN format(' #%s', n) ELSE '' END AS name,
    ROUND((bd.base_price + (random()*80 - 40))::numeric, 2) AS price,
    bd.base_cook + (n % 5) AS cook_time_minutes,
    TRUE,
    'Популярное блюдо категории ' || bd.category_name AS description
FROM base_dishes bd
JOIN cat c ON c.name = bd.category_name
CROSS JOIN LATERAL generate_series(1,3) AS n
ON CONFLICT (category_id, name) DO NOTHING;

----------------------
-- DISH INGREDIENTS --
----------------------
DO $$
DECLARE
    d RECORD;
    prod_ids BIGINT[];
    p_id BIGINT;
    cnt INT;
BEGIN
    FOR d IN SELECT id FROM dishes LOOP
        cnt := 3 + (random()*2)::int; -- от 3 до 5 ингредиентов
        prod_ids := ARRAY(
            SELECT id
            FROM products
            ORDER BY random()
            LIMIT cnt
        );
        FOREACH p_id IN ARRAY prod_ids LOOP
            INSERT INTO dish_ingredients (dish_id, product_id, quantity)
            VALUES (d.id, p_id, ROUND((0.05 + random() * 0.5)::numeric, 3))
            ON CONFLICT (dish_id, product_id) DO NOTHING;
        END LOOP;
    END LOOP;
END;
$$;

------------------
-- RESERVATIONS --
------------------
INSERT INTO reservations (customer_id, table_id, reserved_from, reserved_to, status, created_at)
SELECT
    (SELECT id FROM customers ORDER BY random() LIMIT 1) AS customer_id,
    (SELECT id FROM restaurant_tables ORDER BY random() LIMIT 1) AS table_id,
    rf AS reserved_from,
    rf + interval '2 hours' AS reserved_to,
    (ARRAY['new','confirmed','cancelled','completed'])[1 + floor(random()*4)]::text AS status,
    now() - (random() * interval '10 days') AS created_at
FROM (
    SELECT
        -- день в пределах последних 7 суток
        (date_trunc('day', now())
         - (floor(random()*7))::int * interval '1 day'
         + (12 + floor(random()*10)) * interval '1 hour') AS rf
    FROM generate_series(1,300)
) t
ON CONFLICT DO NOTHING;

------------
-- SHIFTS --
------------
INSERT INTO shifts (opened_by, closed_by, opened_at, closed_at, status, note, expected_revenue, actual_revenue)
SELECT
    (SELECT id FROM employees WHERE role_id = (SELECT id FROM roles WHERE name='manager') ORDER BY random() LIMIT 1),
    (SELECT id FROM employees WHERE role_id = (SELECT id FROM roles WHERE name='manager') ORDER BY random() LIMIT 1),
    (date_trunc('day', now()) - (gs || ' days')::interval) + time '10:00' AS opened_at,
    (date_trunc('day', now()) - (gs || ' days')::interval) + time '23:00' AS closed_at,
    'closed' AS status,
    format('Смена #%s', gs) AS note,
    ROUND((20000 + random()*15000)::numeric, 2) AS expected_revenue,
    ROUND((15000 + random()*20000)::numeric, 2) AS actual_revenue
FROM generate_series(1,60) gs
ON CONFLICT DO NOTHING;

------------
-- ORDERS --
------------
INSERT INTO orders (table_id, customer_id, waiter_id, reservation_id, shift_id, created_at, status)
SELECT
    (SELECT id FROM restaurant_tables ORDER BY random() LIMIT 1) AS table_id,
    (SELECT id FROM customers ORDER BY random() LIMIT 1) AS customer_id,
    (SELECT id FROM employees WHERE role_id = (SELECT id FROM roles WHERE name='waiter') ORDER BY random() LIMIT 1) AS waiter_id,
    NULL::BIGINT AS reservation_id,
    s.id AS shift_id,
    s.opened_at + (s.closed_at - s.opened_at) * random() AS created_at,
    (ARRAY['new','in_progress','closed','cancelled'])[1 + floor(random()*4)] AS status
FROM shifts s,
LATERAL generate_series(1, 1 + (random()*70)::int) gs;

-- Привязываем часть резерваций к заказам (1:1, первых ~200)
WITH ordered_res AS (
    SELECT id, row_number() OVER (ORDER BY created_at) AS rn
    FROM reservations
    LIMIT 200
)
UPDATE orders o
SET reservation_id = r.id
FROM ordered_res r
WHERE o.id = r.rn;

-----------------
-- ORDER ITEMS --
-----------------
INSERT INTO order_items (order_id, dish_id, quantity, price_at_moment, comment)
SELECT
    o.id AS order_id,
    d.id AS dish_id,
    1 + (floor(random()*3))::int AS quantity,
    d.price AS price_at_moment,
    NULL::TEXT AS comment
FROM orders o
CROSS JOIN LATERAL (
    SELECT id, price
    FROM dishes
    ORDER BY random()
    LIMIT (2 + (o.id % 3))         -- 2–4 блюда
) d
ON CONFLICT (order_id, dish_id) DO NOTHING;

--------------
-- PAYMENTS --
--------------
INSERT INTO payments (order_id, amount, method, status, paid_at)
SELECT
    o.id AS order_id,
    COALESCE((
        SELECT SUM(oi.price_at_moment * oi.quantity)
        FROM order_items oi
        WHERE oi.order_id = o.id
    ), 0) AS amount,
    (ARRAY['cash','card','online'])[1 + floor(random()*3)] AS method,
    CASE 
        WHEN o.status = 'cancelled' THEN 'refunded'
        WHEN o.status IN ('new','in_progress') AND random() < 0.5 THEN 'pending'
        ELSE 'paid'
    END AS status,
    o.created_at + interval '30 minutes' AS paid_at
FROM orders o
ON CONFLICT (order_id) DO NOTHING;