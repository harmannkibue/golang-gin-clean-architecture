CREATE TYPE "user_roles" AS ENUM (
  'author',
  'reader'
);

CREATE TABLE "blog" (
                        "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid ()),
                        "descriptions" varchar(15),
                        "user_role" user_roles NOT NULL DEFAULT 'author',
                        "created_at" timestamptz NOT NULL DEFAULT (now()),
                        "updated_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "blog" ("id");
