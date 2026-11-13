CREATE TABLE lti.deep_links (
	id              uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id       uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                            REFERENCES core.tenants (id) ON DELETE CASCADE,
	platform_id     uuid        NOT NULL REFERENCES lti.platforms (id) ON DELETE CASCADE,
	user_id         uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	deployment_id   text        NOT NULL,
	return_url      text        NOT NULL CHECK (return_url ~ '^https?://'),
	data            text        NOT NULL DEFAULT '',
	accept_multiple boolean     NOT NULL DEFAULT true,
	response        text,
	code            text        UNIQUE,
	created_at      timestamptz NOT NULL DEFAULT now()
);

ALTER TABLE lti.deep_links ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.deep_links FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.deep_links
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON lti.deep_links TO lodestar_app;
