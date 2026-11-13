ALTER TABLE study.assignments ADD COLUMN search_map jsonb;

ALTER TABLE study.submissions ADD COLUMN search_map jsonb;

CREATE TABLE study.search_maps (
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	assignment_id uuid        NOT NULL REFERENCES study.assignments (id) ON DELETE CASCADE,
	student_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	data          jsonb       NOT NULL,
	changed_at    timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (assignment_id, student_id)
);

CREATE INDEX search_maps_student ON study.search_maps (tenant_id, student_id);

ALTER TABLE study.search_maps ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.search_maps FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.search_maps
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON study.search_maps TO lodestar_app;
