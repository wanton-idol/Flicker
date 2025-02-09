BEGIN;

# modification in user_search_profile table
ALTER TABLE user_search_profile ADD COLUMN snooze TIMESTAMP DEFAULT NULL;
ALTER TABLE user_search_profile ADD COLUMN hide_my_name boolean;

# modifications in advanced_filters table
ALTER TABLE advanced_filters ADD COLUMN incognito_mode boolean;

COMMIT;