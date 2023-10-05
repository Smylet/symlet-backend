-- Modify "students" table
ALTER TABLE "public"."students" ADD COLUMN "department" text NOT NULL, ADD COLUMN "year_of_entry" timestamptz NULL, ADD COLUMN "expected_graduation_year" timestamptz NULL, ADD COLUMN "student_identification_number" text NOT NULL;
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "profile_id" bigint NULL;
-- Modify "profiles" table
ALTER TABLE "public"."profiles" ADD COLUMN "first_name" text NULL, ADD COLUMN "last_name" text NULL, ADD
 CONSTRAINT "fk_users_profile" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
