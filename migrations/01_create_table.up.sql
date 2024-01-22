BEGIN;

DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "surname" VARCHAR(255) NOT NULL,
    "patronymic" VARCHAR(255),
    "age" INTEGER NOT NULL,
    "gender" VARCHAR(6) NOT NULL,
    "country" VARCHAR(2) NOT NULL,

    CONSTRAINT "user_unique" UNIQUE ("name", "surname", "patronymic", "age", "gender", "country")
);

COMMIT;