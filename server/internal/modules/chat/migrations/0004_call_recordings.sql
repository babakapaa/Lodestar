CREATE TABLE chat.call_recordings (
	egress_id  text        PRIMARY KEY,
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	call_id    uuid        NOT NULL REFERENCES chat.calls (id) ON DELETE CASCADE,
	file_id    uuid        REFERENCES core.files (id) ON DELETE SET NULL,
	started_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	started_at timestamptz NOT NULL DEFAULT clock_timestamp(),
	ended_at   timestamptz
);

CREATE UNIQUE INDEX call_recordings_active ON chat.call_recordings (call_id) WHERE ended_at IS NULL;

ALTER TABLE chat.call_recordings ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.call_recordings FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.call_recordings
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON chat.call_recordings TO lodestar_app;
