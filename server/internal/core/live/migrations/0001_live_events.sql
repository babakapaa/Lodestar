CREATE TABLE core.live_events (
	id        bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	type      text        NOT NULL CHECK (type ~ '^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)+$'),
	payload   jsonb       NOT NULL DEFAULT '{}',
	at        timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX live_events_user ON core.live_events (tenant_id, user_id, id);
CREATE INDEX live_events_at ON core.live_events (tenant_id, at);

ALTER TABLE core.live_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.live_events FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.live_events
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON core.live_events TO lodestar_app;
