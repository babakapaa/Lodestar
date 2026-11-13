CREATE TABLE chat.calls (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	chat_id    uuid        NOT NULL REFERENCES chat.chats (id) ON DELETE CASCADE,
	started_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	started_at timestamptz NOT NULL DEFAULT clock_timestamp(),
	ended_at   timestamptz
);

CREATE UNIQUE INDEX calls_active ON chat.calls (chat_id) WHERE ended_at IS NULL;

CREATE TABLE chat.call_members (
	call_id   uuid        NOT NULL REFERENCES chat.calls (id) ON DELETE CASCADE,
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	joined_at timestamptz NOT NULL DEFAULT clock_timestamp(),
	left_at   timestamptz,
	PRIMARY KEY (call_id, user_id)
);

ALTER TABLE chat.calls ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.calls FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.calls
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE chat.call_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.call_members FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.call_members
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON chat.calls TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON chat.call_members TO lodestar_app;
