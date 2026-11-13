ALTER TABLE chat.members ADD COLUMN pending boolean NOT NULL DEFAULT false;

CREATE INDEX chats_title_search ON chat.chats USING gin (core.search_norm(title) gin_trgm_ops) WHERE kind <> 'direct';

CREATE TABLE chat.message_files (
	tenant_id  uuid NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                REFERENCES core.tenants (id) ON DELETE CASCADE,
	message_id uuid NOT NULL REFERENCES chat.messages (id) ON DELETE CASCADE,
	file_id    uuid NOT NULL REFERENCES core.files (id) ON DELETE CASCADE,
	position   int  NOT NULL,
	PRIMARY KEY (message_id, file_id)
);

CREATE INDEX message_files_file ON chat.message_files (tenant_id, file_id);

CREATE TABLE chat.blocks (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	blocked_id uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, blocked_id),
	CHECK (user_id <> blocked_id)
);

CREATE INDEX blocks_blocked ON chat.blocks (tenant_id, blocked_id);

CREATE TABLE chat.reports (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	chat_id     uuid        REFERENCES chat.chats (id) ON DELETE SET NULL,
	reporter_id uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	data        bytea       NOT NULL,
	resolved_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	resolved_at timestamptz,
	created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX reports_open ON chat.reports (tenant_id, created_at) WHERE resolved_at IS NULL;

ALTER TABLE chat.message_files ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.message_files FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.message_files
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE chat.blocks ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.blocks FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.blocks
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE chat.reports ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.reports FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.reports
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON chat.message_files TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON chat.blocks TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON chat.reports TO lodestar_app;
GRANT DELETE ON chat.messages TO lodestar_app;
