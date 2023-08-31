CREATE TABLE rating_reviews (
    id SERIAL PRIMARY KEY,
    rating NUMERIC(3,2) NOT NULL DEFAULT 0,
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NULL,
    CONSTRAINT fk_redemptions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);