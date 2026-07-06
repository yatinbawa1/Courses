CREATE TYPE "Content_Type" AS ENUM (
  'pdf',
  'quiz',
  'video',
  'external_link',
  'internal_link'
);

CREATE TABLE "User" (
  "user_id" varchar PRIMARY KEY,
  "name" varchar(255),
  "profile_photo_url" text,
  "hashed_password" varchar,
  "revenue" integer
);

CREATE TABLE "Course" (
  "course_id" varchar PRIMARY KEY,
  "course_name" varchar NOT NULL,
  "course_description" varchar NOT NULL,
  "creation_date" Date
);

CREATE TABLE "Course_Tags" (
  "course_id" varchar,
  "tag_id" varchar,          -- Removed the conflicting PRIMARY KEY statement from this line
  "tag" varchar,
  PRIMARY KEY ("course_id", "tag_id")
);

CREATE TABLE "Review" (
  "course_id" varchar,
  "reviewer_id" varchar,
  "review_rating" integer,
  "review_content" varchar,
  PRIMARY KEY ("course_id", "reviewer_id")
);

CREATE TABLE "Course_Chapter" (
  "chapter_id" varchar PRIMARY KEY,
  "course_id" varchar,
  "chapter_index" int,
  "chapter_descrption" varchar,
  "chapter_free" bool,
  "chapter_content" "Content_Type"
);

CREATE TABLE "User_Owned_Courses" (
  "course_id" varchar,
  "user_id" varchar,
  "purchase_date" Date,
  "chapters_consumed_count" int,
  PRIMARY KEY ("course_id", "user_id")
);

CREATE TABLE "User_Chapter_Progress" (
  "user_id" varchar,
  "chapter_id" varchar,
  "is_completed" bool DEFAULT false,
  "consumed_till" integer,
  "last_accessed_at" TIMESTAMP, -- Changed from "DateTime" to native TIMESTAMP
  PRIMARY KEY ("user_id", "chapter_id")
);

CREATE TABLE "User_Created_Courses" (
  "course_id" varchar,
  "user_id" varchar,
  PRIMARY KEY ("course_id", "user_id")
);

CREATE INDEX ON "Course" ("course_name");

ALTER TABLE "Course_Tags" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("course_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "Review" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("course_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "Review" ADD FOREIGN KEY ("reviewer_id") REFERENCES "User" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "Course_Chapter" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("course_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "User_Owned_Courses" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "User_Owned_Courses" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("course_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "User_Chapter_Progress" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "User_Chapter_Progress" ADD FOREIGN KEY ("chapter_id") REFERENCES "Course_Chapter" ("chapter_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "User_Created_Courses" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;
ALTER TABLE "User_Created_Courses" ADD FOREIGN KEY ("course_id") REFERENCES "Course" ("course_id") DEFERRABLE INITIALLY IMMEDIATE;