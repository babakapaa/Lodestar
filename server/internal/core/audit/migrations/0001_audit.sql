CREATE TABLE core.audit_log (
	id          bigint      PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	actor_id    uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	action      text        NOT NULL CHECK (action ~ '^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$'),
	target_kind text        NOT NULL DEFAULT '',
	target_id   text        NOT NULL DEFAULT '',
	details     jsonb       NOT NULL DEFAULT '{}',
	at          timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX audit_log_recent ON core.audit_log (tenant_id, id DESC);
CREATE INDEX audit_log_actor ON core.audit_log (tenant_id, actor_id, id DESC);

ALTER TABLE core.audit_log ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.audit_log FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.audit_log
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON core.audit_log TO lodestar_app;
