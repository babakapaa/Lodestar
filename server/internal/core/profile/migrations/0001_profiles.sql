CREATE TABLE core.profiles (
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id      uuid        PRIMARY KEY REFERENCES core.users (id) ON DELETE CASCADE,
	display_name text        NOT NULL DEFAULT '' CHECK (length(display_name) <= 100),
	about        text        NOT NULL DEFAULT '' CHECK (length(about) <= 2000),
	contacts     jsonb       NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(contacts) = 'array'),
	updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE core.preferences (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	key        text        NOT NULL,
	value      jsonb       NOT NULL,
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (user_id, key)
);

ALTER TABLE core.profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.profiles FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.profiles
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.preferences ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.preferences FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.preferences
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON core.profiles TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON core.preferences TO lodestar_app;
