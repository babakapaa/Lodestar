CREATE TABLE core.launch_tickets (
	ticket_hash bytea       PRIMARY KEY,
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id     uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	device_name text        NOT NULL,
	expires_at  timestamptz NOT NULL
);

ALTER TABLE core.launch_tickets ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.launch_tickets FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.launch_tickets
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON core.launch_tickets TO lodestar_app;
