ALTER TABLE user_nudges DROP COLUMN audio_url;

ALTER TABLE user_nudges ADD COLUMN media_url VARCHAR(255);
ALTER TABLE user_nudges ADD COLUMN type VARCHAR(10);

