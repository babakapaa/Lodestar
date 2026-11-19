ALTER TABLE study.assignments ADD COLUMN navigator_summary boolean NOT NULL DEFAULT false;
ALTER TABLE study.submissions ADD COLUMN navigator_summary jsonb;
