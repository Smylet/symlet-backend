-- Modify "hostel_managers" table
ALTER TABLE "public"."hostel_managers" DROP COLUMN "user_id";
-- Modify "students" table
ALTER TABLE "public"."students" DROP COLUMN "user_id";
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "user_id" bigint NULL, ADD COLUMN "user_type" text NULL;
-- Modify "vendors" table
ALTER TABLE "public"."vendors" DROP COLUMN "user_id";
