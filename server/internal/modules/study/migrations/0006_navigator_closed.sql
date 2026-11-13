ALTER TABLE study.courses DROP CONSTRAINT courses_navigator_mode_check;
ALTER TABLE study.courses ADD CONSTRAINT courses_navigator_mode_check
	CHECK (navigator_mode IN ('study', 'project', 'closed'));

ALTER TABLE study.assignments DROP CONSTRAINT assignments_navigator_mode_check;
ALTER TABLE study.assignments ADD CONSTRAINT assignments_navigator_mode_check
	CHECK (navigator_mode IN ('study', 'project', 'closed'));
