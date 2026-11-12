CREATE TABLE study.series (
	id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	course_id    uuid        NOT NULL REFERENCES study.courses (id) ON DELETE CASCADE,
	unit_id      uuid        REFERENCES core.units (id) ON DELETE CASCADE,
	teacher_id   uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	kind         text        NOT NULL CHECK (kind ~ '^[a-z][a-z0-9_]{0,39}$'),
	room         text        NOT NULL DEFAULT '' CHECK (length(room) <= 200),
	starts_on    date        NOT NULL,
	ends_on      date        NOT NULL,
	weekday      smallint    NOT NULL CHECK (weekday BETWEEN 1 AND 7),
	start_minute smallint    NOT NULL CHECK (start_minute BETWEEN 0 AND 1439),
	minutes      smallint    NOT NULL CHECK (minutes BETWEEN 1 AND 1440),
	every_weeks  smallint    NOT NULL DEFAULT 1 CHECK (every_weeks BETWEEN 1 AND 4),
	created_at   timestamptz NOT NULL DEFAULT now(),
	updated_at   timestamptz NOT NULL DEFAULT now(),
	CHECK (ends_on >= starts_on AND ends_on - starts_on <= 400)
);

CREATE INDEX series_course ON study.series (tenant_id, course_id);
CREATE INDEX series_teacher ON study.series (tenant_id, teacher_id);

CREATE TABLE study.series_changes (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	series_id  uuid        NOT NULL REFERENCES study.series (id) ON DELETE CASCADE,
	on_date    date        NOT NULL,
	status     text        NOT NULL CHECK (status IN ('cancelled', 'changed')),
	starts_at  timestamptz,
	minutes    smallint    CHECK (minutes BETWEEN 1 AND 1440),
	room       text        CHECK (length(room) <= 200),
	teacher_id uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	note       text        NOT NULL DEFAULT '' CHECK (length(note) <= 1000),
	changed_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (series_id, on_date)
);

CREATE TABLE study.reminders (
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	assignment_id uuid        NOT NULL REFERENCES study.assignments (id) ON DELETE CASCADE,
	student_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	sent_at       timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (assignment_id, student_id)
);

CREATE INDEX assignments_due ON study.assignments (tenant_id, due_at) WHERE published_at IS NOT NULL;

ALTER TABLE study.series ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.series FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.series
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.series_changes ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.series_changes FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.series_changes
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.reminders ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.reminders FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.reminders
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON study.series TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON study.series_changes TO lodestar_app;
GRANT SELECT, INSERT ON study.reminders TO lodestar_app;
