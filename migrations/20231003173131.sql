-- Create "sessions" table
CREATE TABLE "public"."sessions" (
  "id" text NOT NULL,
  "username" text NULL,
  "refresh_token" text NULL,
  "user_agent" text NULL,
  "client_ip" text NULL,
  "expires_at" timestamptz NULL,
  "is_blocked" boolean NULL,
  PRIMARY KEY ("id")
);
-- Create "verification_emails" table
CREATE TABLE "public"."verification_emails" (
  "id" bigserial NOT NULL,
  "email" text NULL,
  "secret_code" text NULL,
  "expires_at" timestamptz NULL,
  "user_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create "notifications" table
CREATE TABLE "public"."notifications" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NULL,
  "content" text NULL,
  "is_read" boolean NULL DEFAULT false,
  "action_type" text NULL,
  "action_id" bigint NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_notifications_deleted_at" to table: "notifications"
CREATE INDEX "idx_notifications_deleted_at" ON "public"."notifications" ("deleted_at");
-- Create "profiles" table
CREATE TABLE "public"."profiles" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NULL,
  "bio" text NULL,
  "image" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_profiles_deleted_at" to table: "profiles"
CREATE INDEX "idx_profiles_deleted_at" ON "public"."profiles" ("deleted_at");
-- Create "users" table
CREATE TABLE "public"."users" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "username" text NOT NULL,
  "email" text NOT NULL,
  "password" text NULL,
  "is_email_confirmed" boolean NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX "users_email_key" ON "public"."users" ("email");
-- Create index "users_username_key" to table: "users"
CREATE UNIQUE INDEX "users_username_key" ON "public"."users" ("username");
-- Create "hostel_managers" table
CREATE TABLE "public"."hostel_managers" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostel_managers_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_hostel_managers_deleted_at" to table: "hostel_managers"
CREATE INDEX "idx_hostel_managers_deleted_at" ON "public"."hostel_managers" ("deleted_at");
-- Create "reference_universities" table
CREATE TABLE "public"."reference_universities" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "slug" text NULL,
  "state" text NULL,
  "city" text NULL,
  "country" text NULL,
  "code" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_reference_universities_deleted_at" to table: "reference_universities"
CREATE INDEX "idx_reference_universities_deleted_at" ON "public"."reference_universities" ("deleted_at");
-- Create "hostels" table
CREATE TABLE "public"."hostels" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NOT NULL,
  "university_id" bigint NOT NULL,
  "address" text NOT NULL,
  "city" text NOT NULL,
  "state" text NOT NULL,
  "country" text NOT NULL,
  "description" text NOT NULL,
  "manager_id" bigint NOT NULL,
  "number_of_units" bigint NOT NULL,
  "number_of_occupied_units" bigint NOT NULL,
  "number_of_bedrooms" bigint NOT NULL,
  "number_of_bathrooms" bigint NOT NULL,
  "kitchen" boolean NOT NULL,
  "floor_space" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostels_manager" FOREIGN KEY ("manager_id") REFERENCES "public"."hostel_managers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_hostels_university" FOREIGN KEY ("university_id") REFERENCES "public"."reference_universities" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_hostels_deleted_at" to table: "hostels"
CREATE INDEX "idx_hostels_deleted_at" ON "public"."hostels" ("deleted_at");
-- Create "reference_hostel_ammenities" table
CREATE TABLE "public"."reference_hostel_ammenities" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_reference_hostel_ammenities_deleted_at" to table: "reference_hostel_ammenities"
CREATE INDEX "idx_reference_hostel_ammenities_deleted_at" ON "public"."reference_hostel_ammenities" ("deleted_at");
-- Create "hostel_ammenities" table
CREATE TABLE "public"."hostel_ammenities" (
  "hostel_id" bigint NOT NULL,
  "reference_hostel_ammenities_id" bigint NOT NULL,
  PRIMARY KEY ("hostel_id", "reference_hostel_ammenities_id"),
  CONSTRAINT "fk_hostel_ammenities_hostel" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_hostel_ammenities_reference_hostel_ammenities" FOREIGN KEY ("reference_hostel_ammenities_id") REFERENCES "public"."reference_hostel_ammenities" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "students" table
CREATE TABLE "public"."students" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "university_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_students_university" FOREIGN KEY ("university_id") REFERENCES "public"."reference_universities" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_students_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_students_deleted_at" to table: "students"
CREATE INDEX "idx_students_deleted_at" ON "public"."students" ("deleted_at");
-- Create "hostel_bookings" table
CREATE TABLE "public"."hostel_bookings" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "student_id" bigint NOT NULL,
  "hostel_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostel_bookings_hostel" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_hostel_bookings_student" FOREIGN KEY ("student_id") REFERENCES "public"."students" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_hostel_bookings_deleted_at" to table: "hostel_bookings"
CREATE INDEX "idx_hostel_bookings_deleted_at" ON "public"."hostel_bookings" ("deleted_at");
-- Create "hostel_fees" table
CREATE TABLE "public"."hostel_fees" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "hostel_id" bigint NULL,
  "total_amount" numeric NULL,
  "payment_plan" text NULL,
  "breakdown" json NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostels_hostel_fee" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_hostel_fees_deleted_at" to table: "hostel_fees"
CREATE INDEX "idx_hostel_fees_deleted_at" ON "public"."hostel_fees" ("deleted_at");
-- Create "hostel_images" table
CREATE TABLE "public"."hostel_images" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "image_url" text NOT NULL,
  "hostel_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostels_hostel_images" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_hostel_images_deleted_at" to table: "hostel_images"
CREATE INDEX "idx_hostel_images_deleted_at" ON "public"."hostel_images" ("deleted_at");
-- Create "hostel_maintenance_requests" table
CREATE TABLE "public"."hostel_maintenance_requests" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "hostel_id" bigint NOT NULL,
  "student_id" bigint NOT NULL,
  "resolved_by_id" integer NULL,
  "subject" character varying(64) NOT NULL,
  "description" text NULL,
  "resolve_status" text NULL DEFAULT 'open',
  "resolved" boolean NULL DEFAULT false,
  "resolved_date" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostel_maintenance_requests_hostel" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_hostel_maintenance_requests_resolved_by" FOREIGN KEY ("resolved_by_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_hostel_maintenance_requests_student" FOREIGN KEY ("student_id") REFERENCES "public"."students" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_hostel_maintenance_requests_deleted_at" to table: "hostel_maintenance_requests"
CREATE INDEX "idx_hostel_maintenance_requests_deleted_at" ON "public"."hostel_maintenance_requests" ("deleted_at");
-- Create "hostel_maintenance_request_images" table
CREATE TABLE "public"."hostel_maintenance_request_images" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "image_url" text NOT NULL,
  "hostel_maintenance_request_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostel_maintenance_requests_request_images" FOREIGN KEY ("hostel_maintenance_request_id") REFERENCES "public"."hostel_maintenance_requests" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_hostel_maintenance_request_images_deleted_at" to table: "hostel_maintenance_request_images"
CREATE INDEX "idx_hostel_maintenance_request_images_deleted_at" ON "public"."hostel_maintenance_request_images" ("deleted_at");
-- Create "hostel_manager_reviews" table
CREATE TABLE "public"."hostel_manager_reviews" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "reviewer_id" bigint NOT NULL,
  "manager_id" bigint NOT NULL,
  "rating" numeric NULL DEFAULT 0,
  "description" character varying(1023) NOT NULL,
  "communication_rating" numeric NULL,
  "communication_comment" character varying(127) NULL,
  "professionalism_rating" numeric NULL,
  "professionalism_comment" character varying(127) NULL,
  "responsiveness_rating" numeric NULL,
  "responsiveness_comment" character varying(127) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostel_manager_reviews_manager" FOREIGN KEY ("manager_id") REFERENCES "public"."hostel_managers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_hostel_manager_reviews_reviewer" FOREIGN KEY ("reviewer_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_hostel_manager_reviews_communication_rating" CHECK ((communication_rating >= (1)::numeric) AND (communication_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_manager_reviews_professionalism_rating" CHECK ((professionalism_rating >= (1)::numeric) AND (professionalism_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_manager_reviews_rating" CHECK ((rating >= (0)::numeric) AND (rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_manager_reviews_responsiveness_rating" CHECK ((responsiveness_rating >= (1)::numeric) AND (responsiveness_rating <= (5)::numeric))
);
-- Create index "idx_hostel_manager_reviews_deleted_at" to table: "hostel_manager_reviews"
CREATE INDEX "idx_hostel_manager_reviews_deleted_at" ON "public"."hostel_manager_reviews" ("deleted_at");
-- Create "hostel_reviews" table
CREATE TABLE "public"."hostel_reviews" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "hostel_id" bigint NOT NULL,
  "reviewer_id" bigint NOT NULL,
  "rating" numeric NULL,
  "security_rating" numeric NULL,
  "security_comment" character varying(127) NULL,
  "location_rating" numeric NULL,
  "location_comment" character varying(127) NULL,
  "general_rating" numeric NULL,
  "general_comment" character varying(127) NULL,
  "amenities_rating" numeric NULL,
  "amenities_comment" character varying(127) NULL,
  "water_rating" numeric NULL,
  "water_comment" character varying(127) NULL,
  "electricity_rating" numeric NULL,
  "electricity_comment" character varying(127) NULL,
  "caretaker_rating" numeric NULL,
  "caretaker_comment" character varying(127) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostel_reviews_hostel" FOREIGN KEY ("hostel_id") REFERENCES "public"."hostels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_hostel_reviews_reviewer" FOREIGN KEY ("reviewer_id") REFERENCES "public"."students" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_hostel_reviews_amenities_rating" CHECK ((amenities_rating >= (1)::numeric) AND (amenities_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_reviews_caretaker_rating" CHECK ((caretaker_rating >= (1)::numeric) AND (caretaker_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_reviews_electricity_rating" CHECK ((electricity_rating >= (1)::numeric) AND (electricity_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_reviews_general_rating" CHECK ((general_rating >= (1)::numeric) AND (general_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_reviews_location_rating" CHECK ((location_rating >= (1)::numeric) AND (location_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_reviews_security_rating" CHECK ((security_rating >= (1)::numeric) AND (security_rating <= (5)::numeric)),
  CONSTRAINT "chk_hostel_reviews_water_rating" CHECK ((water_rating >= (1)::numeric) AND (water_rating <= (5)::numeric))
);
-- Create index "idx_hostel_reviews_deleted_at" to table: "hostel_reviews"
CREATE INDEX "idx_hostel_reviews_deleted_at" ON "public"."hostel_reviews" ("deleted_at");
-- Create "hostel_agreement_templates" table
CREATE TABLE "public"."hostel_agreement_templates" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "hostel_id" bigint NOT NULL,
  "document_url" text NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_hostel_agreement_templates_deleted_at" to table: "hostel_agreement_templates"
CREATE INDEX "idx_hostel_agreement_templates_deleted_at" ON "public"."hostel_agreement_templates" ("deleted_at");
-- Create "hostel_students" table
CREATE TABLE "public"."hostel_students" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "student_id" bigint NOT NULL,
  "hostel_id" bigint NOT NULL,
  "check_in_date" timestamptz NULL,
  "check_out_date" timestamptz NULL,
  "room_number" text NULL,
  "current_hostel" boolean NULL,
  "signed_agreement" boolean NULL DEFAULT false,
  "hostel_agreement_template_id" bigint NULL,
  "submitted_signed_agreement_url" text NULL,
  "total_amount_paid" numeric NULL DEFAULT 0,
  "total_amount_due" numeric NULL DEFAULT 0,
  "last_payment_date" timestamptz NULL,
  "next_payment_date" timestamptz NULL,
  "hostel_booking_id" bigint NULL,
  PRIMARY KEY ("id", "student_id", "hostel_id"),
  CONSTRAINT "fk_hostel_students_hostel_agreement_template" FOREIGN KEY ("hostel_agreement_template_id") REFERENCES "public"."hostel_agreement_templates" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_hostel_students_hostel_booking" FOREIGN KEY ("hostel_booking_id") REFERENCES "public"."hostel_bookings" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_hostel_students_deleted_at" to table: "hostel_students"
CREATE INDEX "idx_hostel_students_deleted_at" ON "public"."hostel_students" ("deleted_at");
-- Create "payment_plans" table
CREATE TABLE "public"."payment_plans" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "amount" numeric NOT NULL,
  "hostel_booking_id" bigint NOT NULL,
  "payment_type" text NOT NULL DEFAULT 'all',
  "payment_interval" text NULL,
  "interval_duration" integer NULL,
  "deferred_date" timestamptz NULL,
  "number_of_months" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hostel_bookings_payment_plans" FOREIGN KEY ("hostel_booking_id") REFERENCES "public"."hostel_bookings" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_payment_plans_payment_interval" CHECK (payment_interval = ANY (ARRAY['equal'::text, 'unequal'::text])),
  CONSTRAINT "chk_payment_plans_payment_type" CHECK (payment_type = ANY (ARRAY['all'::text, 'spread'::text, 'stay'::text, 'deferred'::text]))
);
-- Create index "idx_payment_plans_deleted_at" to table: "payment_plans"
CREATE INDEX "idx_payment_plans_deleted_at" ON "public"."payment_plans" ("deleted_at");
-- Create "payment_distributions" table
CREATE TABLE "public"."payment_distributions" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "payment_plan_id" bigint NOT NULL,
  "date" timestamptz NOT NULL,
  "amount" numeric NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_payment_plans_payment_distributions" FOREIGN KEY ("payment_plan_id") REFERENCES "public"."payment_plans" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_payment_distributions_deleted_at" to table: "payment_distributions"
CREATE INDEX "idx_payment_distributions_deleted_at" ON "public"."payment_distributions" ("deleted_at");
-- Create "vendors" table
CREATE TABLE "public"."vendors" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "company_name" text NOT NULL,
  "address" text NOT NULL,
  "email" text NOT NULL,
  "phone" text NOT NULL,
  "website" text NOT NULL,
  "logo" text NOT NULL,
  "description" text NOT NULL,
  "service" text NOT NULL,
  "rating" numeric NOT NULL,
  "is_verified" boolean NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_vendors_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_vendors_deleted_at" to table: "vendors"
CREATE INDEX "idx_vendors_deleted_at" ON "public"."vendors" ("deleted_at");
-- Create "vendor_reviews" table
CREATE TABLE "public"."vendor_reviews" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "vendor_id" bigint NOT NULL,
  "reviewer_id" bigint NOT NULL,
  "rating" numeric NULL DEFAULT 0,
  "general_comment" character varying(127) NULL,
  "general_rating" numeric NULL,
  "quality" numeric NULL,
  "quality_comment" character varying(127) NULL,
  "timeliness" numeric NULL,
  "timeliness_comment" character varying(127) NULL,
  "communication" numeric NULL,
  "communication_comment" character varying(127) NULL,
  "professionalism" numeric NULL,
  "professionalism_comment" character varying(127) NULL,
  "cost_effectiveness" numeric NULL,
  "cost_effectiveness_comment" character varying(127) NULL,
  "reliability" numeric NULL,
  "reliability_comment" character varying(127) NULL,
  "problem_solving" numeric NULL,
  "problem_solving_comment" character varying(127) NULL,
  "flexibility" bigint NULL,
  "flexibility_comment" character varying(127) NULL,
  "customer_satisfaction" numeric NULL,
  "customer_satisfaction_comment" character varying(127) NULL,
  "response_time" numeric NULL,
  "response_time_comment" character varying(127) NULL,
  "problem_resolution" numeric NULL,
  "problem_resolution_comment" character varying(127) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_vendor_reviews_reviewer" FOREIGN KEY ("reviewer_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_vendor_reviews_vendor" FOREIGN KEY ("vendor_id") REFERENCES "public"."vendors" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_vendor_reviews_communication" CHECK ((communication >= (1)::numeric) AND (communication <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_cost_effectiveness" CHECK ((cost_effectiveness >= (1)::numeric) AND (cost_effectiveness <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_customer_satisfaction" CHECK ((customer_satisfaction >= (1)::numeric) AND (customer_satisfaction <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_flexibility" CHECK ((flexibility >= 1) AND (flexibility <= 5)),
  CONSTRAINT "chk_vendor_reviews_problem_resolution" CHECK ((problem_resolution >= (1)::numeric) AND (problem_resolution <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_problem_solving" CHECK ((problem_solving >= (1)::numeric) AND (problem_solving <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_professionalism" CHECK ((professionalism >= (1)::numeric) AND (professionalism <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_quality" CHECK ((quality >= (1)::numeric) AND (quality <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_reliability" CHECK ((reliability >= (1)::numeric) AND (reliability <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_response_time" CHECK ((response_time >= (1)::numeric) AND (response_time <= (5)::numeric)),
  CONSTRAINT "chk_vendor_reviews_timeliness" CHECK ((timeliness >= (1)::numeric) AND (timeliness <= (5)::numeric))
);
-- Create index "idx_vendor_reviews_deleted_at" to table: "vendor_reviews"
CREATE INDEX "idx_vendor_reviews_deleted_at" ON "public"."vendor_reviews" ("deleted_at");
-- Create "work_orders" table
CREATE TABLE "public"."work_orders" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "hostel_maintenance_request_id" bigint NULL,
  "vendor_id" bigint NULL,
  "description" character varying(1023) NULL,
  "status" character varying(16) NULL DEFAULT 'open',
  "cost" numeric NULL DEFAULT 0,
  "completion_date" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_work_orders_hostel_maintenance_request" FOREIGN KEY ("hostel_maintenance_request_id") REFERENCES "public"."hostel_maintenance_requests" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_work_orders_vendor" FOREIGN KEY ("vendor_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_work_orders_deleted_at" to table: "work_orders"
CREATE INDEX "idx_work_orders_deleted_at" ON "public"."work_orders" ("deleted_at");
-- Create "work_order_comments" table
CREATE TABLE "public"."work_order_comments" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "work_order_id" bigint NULL,
  "comment" character varying(1023) NULL,
  "commented_by" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_work_order_comments_commented_by_user" FOREIGN KEY ("commented_by") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_work_orders_comments" FOREIGN KEY ("work_order_id") REFERENCES "public"."work_orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_work_order_comments_deleted_at" to table: "work_order_comments"
CREATE INDEX "idx_work_order_comments_deleted_at" ON "public"."work_order_comments" ("deleted_at");
-- Create "work_order_images" table
CREATE TABLE "public"."work_order_images" (
  "uid" text NULL,
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "image_url" text NOT NULL,
  "work_order_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_work_orders_work_order_images" FOREIGN KEY ("work_order_id") REFERENCES "public"."work_orders" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_work_order_images_deleted_at" to table: "work_order_images"
CREATE INDEX "idx_work_order_images_deleted_at" ON "public"."work_order_images" ("deleted_at");
