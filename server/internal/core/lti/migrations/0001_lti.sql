CREATE SCHEMA IF NOT EXISTS lti;
GRANT USAGE ON SCHEMA lti TO lodestar_app;

CREATE TABLE lti.keys (
	kid         text        PRIMARY KEY,
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	private_key bytea       NOT NULL,
	created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX keys_tenant ON lti.keys (tenant_id, created_at DESC);

CREATE TABLE lti.platforms (
	id             uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id      uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                           REFERENCES core.tenants (id) ON DELETE CASCADE,
	name           text        NOT NULL CHECK (length(name) BETWEEN 1 AND 100),
	issuer         text        NOT NULL CHECK (issuer ~ '^https?://'),
	client_id      text        NOT NULL CHECK (length(client_id) BETWEEN 1 AND 200),
	auth_url       text        NOT NULL CHECK (auth_url ~ '^https?://'),
	token_url      text        NOT NULL CHECK (token_url ~ '^https?://'),
	jwks_url       text        NOT NULL CHECK (jwks_url ~ '^https?://'),
	deployment_ids text[]      NOT NULL DEFAULT '{}',
	match_email    boolean     NOT NULL DEFAULT false,
	created_at     timestamptz NOT NULL DEFAULT now(),
	updated_at     timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, issuer, client_id)
);

CREATE TABLE lti.attempts (
	state       text        PRIMARY KEY,
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	platform_id uuid        NOT NULL REFERENCES lti.platforms (id) ON DELETE CASCADE,
	nonce       text        NOT NULL,
	expires_at  timestamptz NOT NULL
);

CREATE TABLE lti.identities (
	platform_id uuid        NOT NULL REFERENCES lti.platforms (id) ON DELETE CASCADE,
	subject     text        NOT NULL,
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id     uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	created_at  timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (platform_id, subject)
);

CREATE INDEX identities_user ON lti.identities (user_id);

CREATE TABLE lti.contexts (
	id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	platform_id   uuid        NOT NULL REFERENCES lti.platforms (id) ON DELETE CASCADE,
	deployment_id text        NOT NULL,
	context_id    text        NOT NULL,
	title         text        NOT NULL DEFAULT '',
	label         text        NOT NULL DEFAULT '',
	members_url   text,
	lineitems_url text,
	ags_scopes    text[]      NOT NULL DEFAULT '{}',
	origin_kind   text,
	origin_id     text,
	linked_by     uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	linked_at     timestamptz,
	created_at    timestamptz NOT NULL DEFAULT now(),
	updated_at    timestamptz NOT NULL DEFAULT now(),
	UNIQUE (platform_id, deployment_id, context_id),
	CHECK ((origin_kind IS NULL) = (origin_id IS NULL))
);

CREATE INDEX contexts_origin ON lti.contexts (tenant_id, origin_kind, origin_id) WHERE origin_kind IS NOT NULL;

CREATE TABLE lti.members (
	context_id  uuid        NOT NULL REFERENCES lti.contexts (id) ON DELETE CASCADE,
	user_id     uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	roles       text[]      NOT NULL DEFAULT '{}',
	instructor  boolean     NOT NULL DEFAULT false,
	launched_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (context_id, user_id)
);

CREATE INDEX members_user ON lti.members (user_id);

ALTER TABLE lti.keys ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.keys FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.keys
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE lti.platforms ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.platforms FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.platforms
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE lti.attempts ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.attempts FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.attempts
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE lti.identities ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.identities FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.identities
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE lti.contexts ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.contexts FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.contexts
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE lti.members ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.members FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.members
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT ON lti.keys TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON lti.platforms TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON lti.attempts TO lodestar_app;
GRANT SELECT, INSERT ON lti.identities TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON lti.contexts TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON lti.members TO lodestar_app;
