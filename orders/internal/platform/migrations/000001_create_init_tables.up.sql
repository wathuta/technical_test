CREATE TABLE orders (
    order_id VARCHAR(255) PRIMARY KEY,
    customer_id VARCHAR(255),
    items_id JSONB,
    pickup_address JSONB,
    delivery_address JSONB,
    shipping_method VARCHAR(255),
    order_status VARCHAR(255),
    scheduled_pickup_datetime TIMESTAMP,
    scheduled_delivery_datetime TIMESTAMP,
    tracking_number VARCHAR(255),
    payment_method VARCHAR(255),
    invoice_number VARCHAR(255),
    special_instructions TEXT,
    shipping_cost DOUBLE PRECISION,
    insurance_information VARCHAR(255)
);
