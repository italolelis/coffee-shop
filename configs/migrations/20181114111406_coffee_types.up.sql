CREATE TABLE "coffees"
(
    "id"         UUID         NOT NULL,
    "name"       VARCHAR(100) NOT NULL,
    "price"      FLOAT        NOT NULL,
    "created_at" TIMESTAMP    NOT NULL DEFAULT 'now()',
    CONSTRAINT "pk_products" PRIMARY KEY (
                                          "id"
        )
);
