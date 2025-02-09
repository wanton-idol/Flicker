ALTER TABLE profile_media CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

ALTER TABLE profile_media ADD COLUMN image_text TEXT;
ALTER TABLE profile_media ADD COLUMN latitude FLOAT DEFAULT NULL;
ALTER TABLE profile_media ADD COLUMN longitude FLOAT DEFAULT NULL;