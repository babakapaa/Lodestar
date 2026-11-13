CREATE SCHEMA IF NOT EXISTS chat;
GRANT USAGE ON SCHEMA chat TO lodestar_app;

CREATE TABLE chat.chats (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	kind        text        NOT NULL CHECK (kind IN ('direct', 'group', 'channel')),
	title       text        NOT NULL DEFAULT '' CHECK (length(title) <= 200),
	about       text        NOT NULL DEFAULT '' CHECK (length(about) <= 2000),
	origin_kind text        CHECK (length(origin_kind) BETWEEN 1 AND 60),
	origin_id   text        CHECK (length(origin_id) BETWEEN 1 AND 100),
	direct_key  text        CHECK (length(direct_key) BETWEEN 1 AND 100),
	seq         bigint      NOT NULL DEFAULT 0,
	rev         bigint      NOT NULL DEFAULT 0,
	created_by  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	archived_at timestamptz,
	created_at  timestamptz NOT NULL DEFAULT now(),
	changed_at  timestamptz NOT NULL DEFAULT now(),
	CHECK ((origin_kind IS NULL) = (origin_id IS NULL)),
	CHECK ((kind = 'direct') = (direct_key IS NOT NULL)),
	CHECK (kind <> 'direct' OR origin_kind IS NULL),
	UNIQUE (tenant_id, direct_key),
	UNIQUE (tenant_id, origin_kind, origin_id, kind)
);

CREATE TABLE chat.members (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	chat_id   uuid        NOT NULL REFERENCES chat.chats (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	role      text        NOT NULL CHECK (role IN ('owner', 'admin', 'member')),
	synced    boolean     NOT NULL DEFAULT false,
	muted     boolean     NOT NULL DEFAULT false,
	read_seq  bigint      NOT NULL DEFAULT 0,
	joined_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (chat_id, user_id)
);

CREATE INDEX members_user ON chat.members (tenant_id, user_id);
CREATE INDEX members_read ON chat.members (chat_id, read_seq);

CREATE TABLE chat.messages (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	chat_id    uuid        NOT NULL REFERENCES chat.chats (id) ON DELETE CASCADE,
	seq        bigint      NOT NULL,
	rev        bigint      NOT NULL,
	author_id  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	client_id  text        CHECK (length(client_id) BETWEEN 1 AND 100),
	reply_to   bigint,
	data       bytea,
	edited_at  timestamptz,
	deleted_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	CHECK ((data IS NULL) = (deleted_at IS NOT NULL)),
	UNIQUE (chat_id, seq),
	UNIQUE (chat_id, rev),
	UNIQUE (chat_id, author_id, client_id)
);

CREATE INDEX messages_author ON chat.messages (tenant_id, author_id);

ALTER TABLE chat.chats ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.chats FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.chats
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE chat.members ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.members FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.members
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE chat.messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE chat.messages FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON chat.messages
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON chat.chats TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON chat.members TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON chat.messages TO lodestar_app;
