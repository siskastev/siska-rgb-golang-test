CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    category_id INTEGER NOT NULL,
    point INTEGER NOT NULL DEFAULT 0,
    qty INTEGER NOT NULL DEFAULT 0,
    price NUMERIC NOT NULL DEFAULT 0,
    rating NUMERIC NOT NULL default 0,
    descriptions TEXT NOT NULL,
    image VARCHAR(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NULL,
    deleted_at timestamptz NULL,
    CONSTRAINT fk_porducts_category FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_products_point ON products USING btree (point);
CREATE INDEX idx_products_rating ON products USING btree (rating);
CREATE INDEX idx_products_qty ON products USING btree (qty);