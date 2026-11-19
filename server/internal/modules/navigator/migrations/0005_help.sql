CREATE TABLE navigator.help (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	owner_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	scope_kind text        NOT NULL CHECK (length(scope_kind) BETWEEN 1 AND 30),
	scope_id   text        NOT NULL CHECK (length(scope_id) BETWEEN 1 AND 100),
	data       bytea       NOT NULL,
	changed_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, owner_id, scope_kind, scope_id)
);

ALTER TABLE navigator.help ENABLE ROW LEVEL SECURITY;
ALTER TABLE navigator.help FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON navigator.help
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON navigator.help TO lodestar_app;
