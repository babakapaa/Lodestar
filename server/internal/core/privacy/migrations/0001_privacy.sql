CREATE TABLE core.deletion_requests (
	user_id      uuid        PRIMARY KEY REFERENCES core.users (id) ON DELETE CASCADE,
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	requested_at timestamptz NOT NULL DEFAULT now(),
	due_at       timestamptz NOT NULL
);

CREATE INDEX deletion_requests_due ON core.deletion_requests (tenant_id, due_at);

ALTER TABLE core.deletion_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.deletion_requests FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.deletion_requests
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON core.deletion_requests TO lodestar_app;
GRANT DELETE ON core.users TO lodestar_app;
