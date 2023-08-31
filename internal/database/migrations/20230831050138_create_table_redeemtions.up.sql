CREATE TABLE redemptions (
      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      product_id UUID NOT NULL,
      user_id UUID NOT NULL,
      product_name varchar(100) NOT NULL,
      category_name varchar(100) NOT NULL,
      point INTEGER NOT NULL DEFAULT 0,
      descriptions TEXT NOT NULL,
      image VARCHAR(255) NOT NULL,
      created_at timestamptz NOT NULL DEFAULT now(),
      updated_at timestamptz NULL,
      CONSTRAINT fk_redemptions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
      CONSTRAINT unique_product_user UNIQUE (product_id, user_id)
);