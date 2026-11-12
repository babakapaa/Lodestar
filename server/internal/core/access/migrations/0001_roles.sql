CREATE TABLE core.roles (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	name        text        NOT NULL CHECK (length(name) BETWEEN 1 AND 100),
	color       text        NOT NULL DEFAULT '' CHECK (color = '' OR color ~ '^#[0-9A-F]{6}$'),
	icon        text        NOT NULL DEFAULT '' CHECK (length(icon) <= 64),
	position    integer     NOT NULL CHECK (position >= 0),
	permissions text[]      NOT NULL DEFAULT '{}',
	is_default  boolean     NOT NULL DEFAULT false,
	created_at  timestamptz NOT NULL DEFAULT now(),
	updated_at  timestamptz NOT NULL DEFAULT now(),
	CHECK (is_default = (position = 0)),
	CONSTRAINT roles_position UNIQUE (tenant_id, position) DEFERRABLE INITIALLY DEFERRED
);

CREATE UNIQUE INDEX roles_default ON core.roles (tenant_id) WHERE is_default;

CREATE TABLE core.role_assignments (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	role_id    uuid        NOT NULL REFERENCES core.roles (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	scope_kind text        NOT NULL CHECK (scope_kind ~ '^[a-z][a-z0-9_.]{0,39}$'),
	scope_id   text        NOT NULL CHECK (length(scope_id) BETWEEN 1 AND 100),
	expires_at timestamptz,
	granted_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	granted_at timestamptz NOT NULL DEFAULT now(),
	UNIQUE (role_id, user_id, scope_kind, scope_id)
);

CREATE INDEX role_assignments_user ON core.role_assignments (tenant_id, user_id);

CREATE TABLE core.permission_overrides (
	tenant_id  uuid   NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                  REFERENCES core.tenants (id) ON DELETE CASCADE,
	scope_kind text   NOT NULL CHECK (scope_kind ~ '^[a-z][a-z0-9_.]{0,39}$'),
	scope_id   text   NOT NULL CHECK (length(scope_id) BETWEEN 1 AND 100),
	role_id    uuid   NOT NULL REFERENCES core.roles (id) ON DELETE CASCADE,
	allow      text[] NOT NULL DEFAULT '{}',
	deny       text[] NOT NULL DEFAULT '{}',
	PRIMARY KEY (tenant_id, scope_kind, scope_id, role_id)
);

ALTER TABLE core.roles ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.roles FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.roles
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.role_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.role_assignments FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.role_assignments
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.permission_overrides ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.permission_overrides FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.permission_overrides
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON core.roles TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON core.role_assignments TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON core.permission_overrides TO lodestar_app;
