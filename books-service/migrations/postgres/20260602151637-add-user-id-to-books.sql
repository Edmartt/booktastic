
-- +migrate Up
ALTER TABLE books ADD COLUMN user_id VARCHAR(255);
UPDATE books SET user_id = 'unknown' WHERE user_id IS NULL;
ALTER TABLE books ALTER COLUMN user_id SET NOT NULL;

-- +migrate Down
ALTER TABLE books DROP COLUMN user_id;
