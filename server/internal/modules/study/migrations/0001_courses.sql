CREATE SCHEMA IF NOT EXISTS study;
GRANT USAGE ON SCHEMA study TO lodestar_app;

CREATE TABLE study.courses (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	unit_id     uuid        REFERENCES core.units (id) ON DELETE RESTRICT,
	title       text        NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
	code        text        CHECK (length(code) BETWEEN 1 AND 100),
	description text        NOT NULL DEFAULT '' CHECK (length(description) <= 20000),
	archived_at timestamptz,
	created_at  timestamptz NOT NULL DEFAULT now(),
	updated_at  timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, code)
);

CREATE INDEX courses_unit ON study.courses (tenant_id, unit_id);
CREATE INDEX courses_title_search ON study.courses USING gin (core.search_norm(title) gin_trgm_ops) WHERE archived_at IS NULL;
CREATE INDEX courses_code_search ON study.courses USING gin (core.search_norm(coalesce(code, '')) gin_trgm_ops) WHERE archived_at IS NULL;

CREATE TABLE study.course_teachers (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	course_id uuid        NOT NULL REFERENCES study.courses (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	added_at  timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (course_id, user_id)
);

CREATE INDEX course_teachers_user ON study.course_teachers (tenant_id, user_id);

CREATE TABLE study.course_units (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	course_id uuid        NOT NULL REFERENCES study.courses (id) ON DELETE CASCADE,
	unit_id   uuid        NOT NULL REFERENCES core.units (id) ON DELETE CASCADE,
	added_at  timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (course_id, unit_id)
);

CREATE INDEX course_units_unit ON study.course_units (tenant_id, unit_id);

CREATE TABLE study.course_students (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	course_id uuid        NOT NULL REFERENCES study.courses (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	added_at  timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (course_id, user_id)
);

CREATE INDEX course_students_user ON study.course_students (tenant_id, user_id);

ALTER TABLE study.courses ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.courses FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.courses
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.course_teachers ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.course_teachers FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.course_teachers
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.course_units ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.course_units FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.course_units
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.course_students ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.course_students FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.course_students
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON study.courses TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON study.course_teachers TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON study.course_units TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON study.course_students TO lodestar_app;
