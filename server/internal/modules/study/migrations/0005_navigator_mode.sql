ALTER TABLE study.courses ADD COLUMN navigator_mode text NOT NULL DEFAULT 'study'
	CHECK (navigator_mode IN ('study', 'project'));

ALTER TABLE study.assignments ADD COLUMN navigator_mode text
	CHECK (navigator_mode IN ('study', 'project'));
