-- Membuat tabel products
CREATE TABLE IF NOT EXISTS products (
	id VARCHAR(255) NOT NULL DEFAULT gen_random_uuid(),
	user_id VARCHAR(255) NOT NULL,
	name varchar(50) NOT NULL,
	category product_categories NOT NULL,
	qty INTEGER NOT NULL,
	price INTEGER NOT NULL,
	sku	VARCHAR(50) NOT NULL,
	file_id VARCHAR(300) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);