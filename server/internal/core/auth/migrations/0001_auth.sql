CREATE TABLE core.identity_providers (
	id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	key           text        NOT NULL CHECK (key ~ '^[a-z0-9][a-z0-9-]{0,39}$'),
	display_name  text        NOT NULL CHECK (length(display_name) BETWEEN 1 AND 100),
	issuer        text        NOT NULL,
	client_id     text        NOT NULL,
	client_secret bytea       NOT NULL,
	scopes        text[]      NOT NULL DEFAULT ARRAY['openid', 'email', 'profile'],
	created_at    timestamptz NOT NULL DEFAULT now(),
	updated_at    timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, key)
);

CREATE TABLE core.login_attempts (
	state          text        PRIMARY KEY,
	tenant_id      uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                           REFERENCES core.tenants (id) ON DELETE CASCADE,
	provider_id    uuid        NOT NULL REFERENCES core.identity_providers (id) ON DELETE CASCADE,
	nonce          text        NOT NULL,
	verifier       text        NOT NULL,
	redirect_uri   text        NOT NULL,
	code_challenge text        NOT NULL,
	device_name    text        NOT NULL,
	expires_at     timestamptz NOT NULL
);

CREATE TABLE core.login_grants (
	code_hash      bytea       PRIMARY KEY,
	tenant_id      uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                           REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id        uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	redirect_uri   text        NOT NULL,
	code_challenge text        NOT NULL,
	device_name    text        NOT NULL,
	expires_at     timestamptz NOT NULL
);

CREATE TABLE core.sessions (
	id                    uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id             uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                                  REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id               uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	device_name           text        NOT NULL,
	access_hash           bytea       NOT NULL UNIQUE,
	access_expires_at     timestamptz NOT NULL,
	refresh_hash          bytea       NOT NULL UNIQUE,
	previous_refresh_hash bytea       UNIQUE,
	refresh_expires_at    timestamptz NOT NULL,
	created_at            timestamptz NOT NULL DEFAULT now(),
	last_used_at          timestamptz NOT NULL DEFAULT now(),
	revoked_at            timestamptz
);

CREATE INDEX sessions_user ON core.sessions (tenant_id, user_id) WHERE revoked_at IS NULL;

ALTER TABLE core.identity_providers ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.identity_providers FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.identity_providers
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.login_attempts ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.login_attempts FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.login_attempts
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.login_grants ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.login_grants FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.login_grants
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.sessions FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.sessions
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON core.identity_providers TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON core.login_attempts TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON core.login_grants TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON core.sessions TO lodestar_app;
