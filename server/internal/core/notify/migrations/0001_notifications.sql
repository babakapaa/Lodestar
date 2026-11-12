CREATE TABLE core.notifications (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	type       text        NOT NULL,
	params     jsonb       NOT NULL DEFAULT '{}',
	link       text        NOT NULL DEFAULT '' CHECK (length(link) <= 500),
	created_at timestamptz NOT NULL DEFAULT now(),
	read_at    timestamptz
);

CREATE INDEX notifications_user ON core.notifications (tenant_id, user_id, created_at DESC, id);
CREATE INDEX notifications_unread ON core.notifications (tenant_id, user_id) WHERE read_at IS NULL;

CREATE TABLE core.notification_channels (
	tenant_id uuid    NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                  REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id   uuid    NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	type      text    NOT NULL,
	in_app    boolean NOT NULL,
	push      boolean NOT NULL,
	email     boolean NOT NULL,
	PRIMARY KEY (user_id, type)
);

ALTER TABLE core.notifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.notifications FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.notifications
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.notification_channels ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.notification_channels FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.notification_channels
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON core.notifications TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON core.notification_channels TO lodestar_app;
