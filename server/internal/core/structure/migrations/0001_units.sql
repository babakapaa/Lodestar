CREATE TABLE core.units (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	parent_id   uuid        REFERENCES core.units (id) ON DELETE RESTRICT,
	kind        text        NOT NULL CHECK (kind ~ '^[a-z][a-z0-9_]{0,39}$'),
	name        text        NOT NULL CHECK (length(name) BETWEEN 1 AND 200),
	code        text        CHECK (length(code) BETWEEN 1 AND 100),
	archived_at timestamptz,
	created_at  timestamptz NOT NULL DEFAULT now(),
	updated_at  timestamptz NOT NULL DEFAULT now(),
	CHECK (parent_id IS DISTINCT FROM id),
	UNIQUE (tenant_id, code)
);

CREATE INDEX units_parent ON core.units (tenant_id, parent_id);

CREATE TABLE core.unit_members (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	unit_id   uuid        NOT NULL REFERENCES core.units (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	joined_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (unit_id, user_id)
);

CREATE INDEX unit_members_user ON core.unit_members (tenant_id, user_id);

CREATE TABLE core.unit_events (
	id        bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	unit_id   uuid        NOT NULL REFERENCES core.units (id) ON DELETE CASCADE,
	actor_id  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	action    text        NOT NULL,
	data      jsonb       NOT NULL DEFAULT '{}',
	at        timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX unit_events_unit ON core.unit_events (tenant_id, unit_id, id);

ALTER TABLE core.units ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.units FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.units
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.unit_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.unit_members FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.unit_members
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.unit_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.unit_events FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.unit_events
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON core.units TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON core.unit_members TO lodestar_app;
GRANT SELECT, INSERT ON core.unit_events TO lodestar_app;
