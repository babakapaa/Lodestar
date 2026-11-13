CREATE SCHEMA IF NOT EXISTS attendance;
GRANT USAGE ON SCHEMA attendance TO lodestar_app;

CREATE TABLE attendance.sessions (
	id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	series_id    text        NOT NULL CHECK (length(series_id) BETWEEN 1 AND 100),
	lesson_date  date        NOT NULL,
	course_id    text        NOT NULL CHECK (length(course_id) BETWEEN 1 AND 100),
	course_title text        NOT NULL DEFAULT '' CHECK (length(course_title) <= 300),
	room         text        NOT NULL DEFAULT '' CHECK (length(room) <= 100),
	starts_at    timestamptz NOT NULL,
	ends_at      timestamptz NOT NULL,
	secret       bytea       NOT NULL,
	code         text        NOT NULL CHECK (code ~ '^[0-9]{4}$'),
	code_until   timestamptz NOT NULL,
	opened_by    uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	opened_at    timestamptz NOT NULL DEFAULT clock_timestamp(),
	closed_at    timestamptz,
	CHECK (ends_at > starts_at),
	UNIQUE (tenant_id, series_id, lesson_date)
);

CREATE INDEX sessions_open ON attendance.sessions (tenant_id, code) WHERE closed_at IS NULL;
CREATE INDEX sessions_course ON attendance.sessions (tenant_id, course_id, lesson_date);

CREATE TABLE attendance.marks (
	session_id uuid        NOT NULL REFERENCES attendance.sessions (id) ON DELETE CASCADE,
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	status     text        CHECK (status IN ('present', 'late', 'absent', 'excused')),
	method     text        CHECK (method IN ('qr', 'code', 'list', 'manual')),
	marked_by  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	marked_at  timestamptz,
	PRIMARY KEY (session_id, user_id),
	CHECK ((status IS NULL) = (method IS NULL) AND (status IS NULL) = (marked_at IS NULL))
);

CREATE INDEX marks_user ON attendance.marks (tenant_id, user_id);

CREATE TABLE attendance.failures (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	failed_at timestamptz NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX failures_user ON attendance.failures (tenant_id, user_id, failed_at);

ALTER TABLE attendance.sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE attendance.sessions FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON attendance.sessions
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE attendance.marks ENABLE ROW LEVEL SECURITY;
ALTER TABLE attendance.marks FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON attendance.marks
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE attendance.failures ENABLE ROW LEVEL SECURITY;
ALTER TABLE attendance.failures FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON attendance.failures
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON attendance.sessions TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON attendance.marks TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON attendance.failures TO lodestar_app;
