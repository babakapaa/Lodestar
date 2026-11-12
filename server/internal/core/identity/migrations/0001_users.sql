CREATE TABLE core.users (
	id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	issuer       text        NOT NULL,
	subject      text        NOT NULL,
	email        text        NOT NULL DEFAULT '',
	display_name text        NOT NULL DEFAULT '',
	created_at   timestamptz NOT NULL DEFAULT now(),
	updated_at   timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, issuer, subject)
);

ALTER TABLE core.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.users FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.users
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON core.users TO lodestar_app;
