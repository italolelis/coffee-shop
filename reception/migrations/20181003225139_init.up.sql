CREATE TABLE "orders" (
    "id" UUID NOT NULL,
    "items" JSON NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT 'now()',
    CONSTRAINT "pk_orders" PRIMARY KEY (
        "id"
     )
);
