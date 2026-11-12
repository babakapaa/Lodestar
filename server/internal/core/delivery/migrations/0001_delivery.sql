CREATE TABLE core.push_devices (
	id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id      uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	platform     text        NOT NULL CHECK (platform IN ('android', 'ios', 'web')),
	token        text        NOT NULL CHECK (length(token) BETWEEN 1 AND 4096),
	name         text        NOT NULL DEFAULT '' CHECK (length(name) <= 100),
	language     text        NOT NULL DEFAULT '' CHECK (length(language) <= 10),
	created_at   timestamptz NOT NULL DEFAULT now(),
	last_seen_at timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, token)
);

CREATE INDEX push_devices_user ON core.push_devices (tenant_id, user_id);

CREATE TABLE core.deliveries (
	id         bigint      PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	channel    text        NOT NULL CHECK (channel IN ('push', 'email')),
	type       text        NOT NULL,
	params     jsonb       NOT NULL DEFAULT '{}',
	link       text        NOT NULL DEFAULT '',
	attempts   integer     NOT NULL DEFAULT 0,
	next_at    timestamptz NOT NULL DEFAULT now(),
	last_error text        NOT NULL DEFAULT '',
	created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX deliveries_due ON core.deliveries (tenant_id, next_at);

ALTER TABLE core.push_devices ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.push_devices FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.push_devices
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.deliveries ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.deliveries FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.deliveries
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON core.push_devices TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON core.deliveries TO lodestar_app;
