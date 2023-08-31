CREATE TABLE product_categories (
   id SERIAL PRIMARY KEY,
   name VARCHAR(100) NOT NULL,
   created_at timestamptz NOT NULL DEFAULT now(),
   updated_at timestamptz NULL
);

CREATE INDEX idx_categories_name ON product_categories USING btree (name);