ALTER TABLE "product" ADD CONSTRAINT check_stock_positive CHECK (stock >= 0);
