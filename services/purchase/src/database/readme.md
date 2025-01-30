# Readme first
`migration_generator.go` adalah program pembantu untuk generate `.sql`

# Usage
tinggal 

```bash
# purchase(id, sender_name, sender_contact_detail, sender_contact_type) 
go run migration_generator.go purchase "id:VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),sender_name:VARCHAR(255),sender_contact_detail:VARCHAR(255),sender_contact_type:VARCHAR(10)"

# purchase_cart (purchase_id, product_id, quantity) 
go run migration_generator.go purchase_cart "purchase_id:VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),product_id:VARCHAR(255),quantity:INT"

# kalau mau multiple primary key (yang satunya ber FK)
# go run migration_generator.go purchase_cart "purchase_id:INT,product_id:TEXT,quantity:INT,PRIMARY KEY (purchase_id, product_id)"
```