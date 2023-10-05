-- Create "reference_hostel_amenities" table
CREATE TABLE "public"."reference_hostel_amenities" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_reference_hostel_amenities_deleted_at" to table: "reference_hostel_amenities"
CREATE INDEX "idx_reference_hostel_amenities_deleted_at" ON "public"."reference_hostel_amenities" ("deleted_at");
-- Modify "hostel_ammenities" table
ALTER TABLE "public"."hostel_ammenities" DROP CONSTRAINT "hostel_ammenities_pkey", DROP COLUMN "reference_hostel_ammenities_id", ADD COLUMN "reference_hostel_amenities_id" bigint NOT NULL, ADD PRIMARY KEY ("hostel_id", "reference_hostel_amenities_id"), ADD
 CONSTRAINT "fk_hostel_ammenities_reference_hostel_amenities" FOREIGN KEY ("reference_hostel_amenities_id") REFERENCES "public"."reference_hostel_amenities" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "hostels" table
ALTER TABLE "public"."hostels" DROP CONSTRAINT "fk_hostels_manager", DROP CONSTRAINT "fk_hostels_university", ALTER COLUMN "kitchen" TYPE text, ADD
 CONSTRAINT "fk_hostels_manager" FOREIGN KEY ("manager_id") REFERENCES "public"."hostel_managers" ("id") ON UPDATE CASCADE ON DELETE CASCADE, ADD
 CONSTRAINT "fk_hostels_university" FOREIGN KEY ("university_id") REFERENCES "public"."reference_universities" ("id") ON UPDATE CASCADE ON DELETE SET NULL;
-- Modify "hostel_fees" table
ALTER TABLE "public"."hostel_fees" DROP CONSTRAINT "fk_hostels_hostel_fee", ALTER COLUMN "breakdown" TYPE jsonb, ADD
 CONSTRAINT "fk_hostels_hostel_fee" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "hostel_images" table
ALTER TABLE "public"."hostel_images" DROP CONSTRAINT "fk_hostels_hostel_images", ADD
 CONSTRAINT "fk_hostels_hostel_images" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Drop "reference_hostel_ammenities" table
DROP TABLE "public"."reference_hostel_ammenities";
