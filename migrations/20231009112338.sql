-- Modify "users" table
ALTER TABLE "public"."users" DROP CONSTRAINT "fk_users_profile", ADD
 CONSTRAINT "fk_users_profile" FOREIGN KEY ("profile_id") REFERENCES "public"."profiles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
