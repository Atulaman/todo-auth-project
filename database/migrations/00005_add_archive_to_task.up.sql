BEGIN;
ALTER TABLE "Tasks" ADD COLUMN archive BOOLEAN DEFAULT FALSE;
COMMIT;