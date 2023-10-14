-- Modify "users" table
ALTER TABLE "public"."users" DROP CONSTRAINT "fk_users_profile";
-- Modify "profiles" table
ALTER TABLE "public"."profiles" ADD
 CONSTRAINT "fk_users_profile" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
