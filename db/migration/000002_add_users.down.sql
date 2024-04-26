DROP TABLE IF EXISTS users;


ALTER TABLE IF EXISTS "accounts" DROP FOREIGN KEY ("owner");

ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT "owner_currency_key";