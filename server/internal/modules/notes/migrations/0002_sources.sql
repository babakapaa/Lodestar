CREATE TABLE notes.sources (
	id         uuid        PRIMARY KEY,
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	owner_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	data       bytea       NOT NULL,
	file_id    uuid        REFERENCES core.files (id) ON DELETE SET NULL,
	bytes      bigint      NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	changed_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX sources_owner ON notes.sources (tenant_id, owner_id);

ALTER TABLE notes.sources ENABLE ROW LEVEL SECURITY;
ALTER TABLE notes.sources FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON notes.sources
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);
GRANT SELECT, INSERT, UPDATE, DELETE ON notes.sources TO lodestar_app;
