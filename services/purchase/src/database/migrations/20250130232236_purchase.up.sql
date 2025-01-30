CREATE TABLE purchase (
    id VARCHAR(255) PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    sender_name VARCHAR(255),
    sender_contact_detail VARCHAR(255),
    sender_contact_type VARCHAR(10)
);
