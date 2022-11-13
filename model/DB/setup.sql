CREATE TABLE IF NOT EXISTS "tb_user" (
    "id" SERIAL PRIMARY KEY,
    "username" TEXT NOT NULL UNIQUE,
    "display_name" TEXT,
    "profile" TEXT,
    "avatar_url" TEXT,
    "created_at" TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "tb_account" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL REFERENCES "tb_user" (id),
    "password" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "tb_email" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INTEGER NOT NULL REFERENCES "tb_user" (id),
    "email" TEXT NOT NULL UNIQUE,
    "is_primary" BOOLEAN NOT NULL,
    "is_verified" BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS "tb_dictionary" (
    "id" SERIAL PRIMARY KEY,
    "dictionary_name" TEXT NOT NULL UNIQUE,
    "dictionary_display_name" TEXT,
    "owner_id" INTEGER NOT NULL REFERENCES "tb_user" (id),
    "description" TEXT,
    "image_url" TEXT,
    "scansion_url" TEXT,
    "created_at" TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "tb_word" (
    "id" SERIAL PRIMARY KEY,
    "dictionary_id" INTEGER NOT NULL REFERENCES "tb_dictionary" (id),
    "added_by" INTEGER NOT NULL REFERENCES "tb_user" (id),
    "created_at" TIMESTAMP,
    "updated_at" TIMESTAMP,
    "deleted_at" TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "tb_word_snapshot" (
    "id" SERIAL PRIMARY KEY,
    "word_id" INTEGER NOT NULL REFERENCES "tb_word" (id),
    "headword" text NOT NULL,
    "translation" TEXT,
    "example" TEXT,
    "edited_by" integer not null references "tb_user" (id),
    "updated_at" TIMESTAMP
);

