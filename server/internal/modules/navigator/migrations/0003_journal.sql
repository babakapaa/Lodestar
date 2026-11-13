CREATE TABLE navigator.journal (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	owner_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	scope_kind text        NOT NULL DEFAULT '',
	scope_id   text        NOT NULL DEFAULT '',
	data       bytea       NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX journal_owner ON navigator.journal (tenant_id, owner_id, scope_kind, scope_id, created_at);

ALTER TABLE navigator.journal ENABLE ROW LEVEL SECURITY;
ALTER TABLE navigator.journal FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON navigator.journal
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON navigator.journal TO lodestar_app;
