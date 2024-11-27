BEGIN;
ALTER TABLE session RENAME COLUMN create_at TO created_at;
COMMIT;