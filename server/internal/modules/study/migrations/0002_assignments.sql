CREATE TABLE study.assignments (
	id           uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id    uuid         NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	course_id    uuid         NOT NULL REFERENCES study.courses (id) ON DELETE CASCADE,
	title        text         NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
	description  text         NOT NULL DEFAULT '' CHECK (length(description) <= 50000),
	due_at       timestamptz,
	grading      text         NOT NULL CHECK (grading IN ('points', 'pass', 'none')),
	max_points   numeric(9,2) CHECK (max_points > 0),
	published_at timestamptz,
	created_by   uuid         REFERENCES core.users (id) ON DELETE SET NULL,
	created_at   timestamptz  NOT NULL DEFAULT now(),
	updated_at   timestamptz  NOT NULL DEFAULT now(),
	CHECK ((grading = 'points') = (max_points IS NOT NULL))
);

CREATE INDEX assignments_course ON study.assignments (tenant_id, course_id, due_at);

CREATE TABLE study.assignment_files (
	tenant_id     uuid NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                   REFERENCES core.tenants (id) ON DELETE CASCADE,
	assignment_id uuid NOT NULL REFERENCES study.assignments (id) ON DELETE CASCADE,
	file_id       uuid NOT NULL REFERENCES core.files (id) ON DELETE CASCADE,
	position      int  NOT NULL,
	PRIMARY KEY (assignment_id, file_id)
);

CREATE INDEX assignment_files_file ON study.assignment_files (tenant_id, file_id);

CREATE TABLE study.submissions (
	id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	assignment_id uuid        NOT NULL REFERENCES study.assignments (id) ON DELETE CASCADE,
	student_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	client_id     uuid        NOT NULL,
	body          text        NOT NULL DEFAULT '' CHECK (length(body) <= 100000),
	composed_at   timestamptz,
	submitted_at  timestamptz NOT NULL DEFAULT now(),
	UNIQUE (assignment_id, student_id, client_id)
);

CREATE INDEX submissions_student ON study.submissions (tenant_id, assignment_id, student_id, submitted_at DESC);

CREATE TABLE study.submission_files (
	tenant_id     uuid NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                   REFERENCES core.tenants (id) ON DELETE CASCADE,
	submission_id uuid NOT NULL REFERENCES study.submissions (id) ON DELETE CASCADE,
	file_id       uuid NOT NULL REFERENCES core.files (id) ON DELETE CASCADE,
	position      int  NOT NULL,
	PRIMARY KEY (submission_id, file_id)
);

CREATE INDEX submission_files_file ON study.submission_files (tenant_id, file_id);

CREATE TABLE study.grades (
	tenant_id     uuid         NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                           REFERENCES core.tenants (id) ON DELETE CASCADE,
	assignment_id uuid         NOT NULL REFERENCES study.assignments (id) ON DELETE CASCADE,
	student_id    uuid         NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	points        numeric(9,2) CHECK (points >= 0),
	passed        boolean,
	comment       text         NOT NULL DEFAULT '' CHECK (length(comment) <= 20000),
	graded_by     uuid         REFERENCES core.users (id) ON DELETE SET NULL,
	graded_at     timestamptz  NOT NULL DEFAULT now(),
	PRIMARY KEY (assignment_id, student_id)
);

CREATE INDEX grades_student ON study.grades (tenant_id, student_id);

ALTER TABLE study.assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.assignments FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.assignments
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.assignment_files ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.assignment_files FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.assignment_files
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.submissions ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.submissions FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.submissions
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.submission_files ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.submission_files FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.submission_files
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.grades ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.grades FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.grades
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON study.assignments TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON study.assignment_files TO lodestar_app;
GRANT SELECT, INSERT ON study.submissions TO lodestar_app;
GRANT SELECT, INSERT ON study.submission_files TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON study.grades TO lodestar_app;
