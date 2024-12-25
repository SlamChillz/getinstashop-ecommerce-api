ALTER TABLE "products" ADD CONSTRAINT check_price_positive CHECK (price > 0.0);
