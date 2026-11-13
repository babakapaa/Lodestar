CREATE SCHEMA IF NOT EXISTS navigator;
GRANT USAGE ON SCHEMA navigator TO lodestar_app;

CREATE TABLE navigator.dialogs (
	id         uuid        PRIMARY KEY,
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	owner_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	scope_kind text        NOT NULL DEFAULT '' CHECK (length(scope_kind) <= 30),
	scope_id   text        NOT NULL DEFAULT '' CHECK (length(scope_id) <= 100),
	title      bytea,
	created_at timestamptz NOT NULL DEFAULT now(),
	changed_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX dialogs_owner ON navigator.dialogs (tenant_id, owner_id, scope_kind, scope_id, changed_at DESC);

CREATE TABLE navigator.messages (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	dialog_id  uuid        NOT NULL REFERENCES navigator.dialogs (id) ON DELETE CASCADE,
	seq        int         NOT NULL,
	client_key text        CHECK (length(client_key) BETWEEN 1 AND 100),
	role       text        NOT NULL CHECK (role IN ('student', 'navigator')),
	data       bytea       NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (dialog_id, seq),
	UNIQUE (dialog_id, client_key)
);

ALTER TABLE navigator.dialogs ENABLE ROW LEVEL SECURITY;
ALTER TABLE navigator.dialogs FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON navigator.dialogs
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE navigator.messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE navigator.messages FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON navigator.messages
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON navigator.dialogs TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON navigator.messages TO lodestar_app;
