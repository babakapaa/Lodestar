CREATE TABLE core.tenant_settings (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	key        text        NOT NULL,
	value      jsonb       NOT NULL,
	updated_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, key)
);

CREATE TABLE core.tenant_modules (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	module     text        NOT NULL,
	enabled    boolean     NOT NULL,
	updated_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, module)
);

ALTER TABLE core.tenant_settings ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.tenant_settings FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.tenant_settings
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.tenant_modules ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.tenant_modules FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.tenant_modules
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON core.tenant_settings TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON core.tenant_modules TO lodestar_app;
