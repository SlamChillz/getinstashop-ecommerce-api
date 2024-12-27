CREATE TABLE "user" (
   "id" UUID PRIMARY KEY,  -- Unique identifier for the user
   "email" VARCHAR(255) UNIQUE NOT NULL,  -- User's email, unique and cannot be null
   "password" VARCHAR(255) NOT NULL,  -- User's password, cannot be null
   "admin" BOOLEAN NOT NULL DEFAULT FALSE,  -- Whether the user is an admin, defaults to FALSE
   "createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),  -- Timestamp of when the user was created, defaults to the current time
   "updatedAt" TIMESTAMP NOT NULL DEFAULT NOW()  -- Timestamp of when the user was last updated, defaults to the current time
);

CREATE TABLE "product" (
    "id" UUID PRIMARY KEY,  -- Unique identifier for the product
    "name" VARCHAR(255) UNIQUE NOT NULL,  -- Product name, unique and cannot be null
    "description" TEXT NOT NULL,  -- Product description, cannot be null
    "price" FLOAT NOT NULL,  -- Price of the product, cannot be null
    "stock" INT NOT NULL,  -- Available stock for the product, cannot be null
    "createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),  -- Timestamp of when the product was created, defaults to current time
    "updatedAt" TIMESTAMP NOT NULL DEFAULT NOW(),  -- Timestamp of when the product was last updated, defaults to current time
    "createdBy" UUID NOT NULL,  -- UUID of the user who created the product, cannot be null
    CONSTRAINT "fk_user" FOREIGN KEY ("createdBy") REFERENCES "user"("id")  -- Foreign key relationship with user table
      ON DELETE RESTRICT  -- Ensures that deleting the user(admin) that created the product is restricted.
);

CREATE TYPE "order_status" AS ENUM ('PENDING', 'COMPLETED', 'CANCELLED');

CREATE TABLE "order" (
    "id" UUID PRIMARY KEY,  -- Unique identifier for the order
    "userId" UUID NOT NULL,  -- UUID of the user who placed the order
    "total" FLOAT NOT NULL,  -- Total amount for the order
    "status" "order_status" NOT NULL,  -- Order status (PENDING, COMPLETED, CANCELLED)
    "createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),  -- Timestamp of when the order was created
    "updatedAt" TIMESTAMP NOT NULL DEFAULT NOW(),  -- Timestamp of when the order was last updated
    CONSTRAINT "fk_user" FOREIGN KEY ("userId") REFERENCES "user"("id")  -- Foreign key referencing the user table
        ON DELETE CASCADE  -- Ensures that order are deleted if the associated user is deleted
);

CREATE TABLE "orderItem" (
    "id" UUID PRIMARY KEY,  -- Unique identifier for the order item
    "orderId" UUID NOT NULL,  -- UUID of the order to which this item belongs
    "productId" UUID NOT NULL,  -- UUID of the product being ordered
    "quantity" INT NOT NULL,  -- Quantity of the product ordered
    "price" FLOAT NOT NULL,  -- Price of the product for this order item
    "createdAt" TIMESTAMP NOT NULL DEFAULT NOW(),  -- Timestamp of when the order item was created
    "updatedAt" TIMESTAMP NOT NULL DEFAULT NOW(),  -- Timestamp of when the order item was last updated
    CONSTRAINT "fk_order" FOREIGN KEY ("orderId") REFERENCES "order"("id")  -- Foreign key referencing the order table
       ON DELETE CASCADE,  -- Ensures that order items are deleted if the associated order is deleted
    CONSTRAINT "fk_product" FOREIGN KEY ("productId") REFERENCES "product"("id")  -- Foreign key referencing the product table
       ON DELETE CASCADE  -- Ensures that order items are deleted if the associated product is deleted
);

CREATE OR REPLACE FUNCTION check_admin() RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT "admin" FROM "user" WHERE id = NEW."createdBy") = FALSE THEN
        RAISE EXCEPTION 'Only an admin can create or update a product';
    END IF;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_admin_before_product_write
BEFORE INSERT OR UPDATE ON "product"
FOR EACH ROW
EXECUTE FUNCTION check_admin();
