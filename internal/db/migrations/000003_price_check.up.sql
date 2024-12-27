ALTER TABLE "product" ADD CONSTRAINT check_price_positive CHECK (price > 0.0);
