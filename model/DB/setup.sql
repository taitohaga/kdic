CREATE TABLE IF NOT EXISTS "tb_user" (
    "id" INTEGER AUTO_INCREMENT PRIMARY KEY,
    "username" TEXT NOT NULL,
    "display_name" TEXT,
    "profile" TEXT
);

CREATE TABLE IF NOT EXISTS "tb_account" (
    "id" INTEGER AUTO_INCREMENT PRIMARY KEY,
    "user_id" INTEGER NOT NULL references "tb_user" (id),
    "password" TEXT
);

create table if not exists "tb_email" (
    "id" INTEGER AUTO_INCREMENT PRIMARY KEY,
    "user_id" INTEGER NOT NULL references "tb_user" (id),
    "email" TEXT NOT NULL UNIQUE
);
