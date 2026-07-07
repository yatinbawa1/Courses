-- 000001_init.down.sql

-- Drop tables with CASCADE to automatically remove all foreign key constraints safely
DROP TABLE IF EXISTS "User_Created_Courses" CASCADE;
DROP TABLE IF EXISTS "User_Chapter_Progress" CASCADE;
DROP TABLE IF EXISTS "User_Owned_Courses" CASCADE;
DROP TABLE IF EXISTS "Review" CASCADE;
DROP TABLE IF EXISTS "Course_Tags" CASCADE;
DROP TABLE IF EXISTS "Course_Chapter" CASCADE;
DROP TABLE IF EXISTS "Course" CASCADE;
DROP TABLE IF EXISTS "User" CASCADE;

-- Drop Custom Enums and Types
DROP TYPE IF EXISTS "Content_Type" CASCADE;