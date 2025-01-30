CREATE TABLE purchase_cart (
    purchase_id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    product_id VARCHAR(255),
    quantity INT
);
