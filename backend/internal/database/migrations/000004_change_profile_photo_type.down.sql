ALTER TABLE "User"
    DROP COLUMN IF EXISTS profile_photo_exists,
    ADD COLUMN profile_photo_url TEXT;