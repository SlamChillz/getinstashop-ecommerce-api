ALTER TABLE "products" ADD CONSTRAINT check_stock_positive CHECK (stock >= 0);
