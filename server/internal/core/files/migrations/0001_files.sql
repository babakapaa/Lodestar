CREATE TABLE core.files (
	id           uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	owner_id     uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	name         text        NOT NULL CHECK (length(name) BETWEEN 1 AND 255),
	content_type text        NOT NULL CHECK (length(content_type) BETWEEN 1 AND 255),
	size         bigint      NOT NULL CHECK (size > 0),
	status       text        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'ready')),
	created_at   timestamptz NOT NULL DEFAULT now(),
	uploaded_at  timestamptz
);

CREATE INDEX files_owner ON core.files (tenant_id, owner_id, created_at DESC, id) WHERE status = 'ready';
CREATE INDEX files_pending ON core.files (tenant_id, created_at) WHERE status = 'pending';

ALTER TABLE core.files ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.files FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.files
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON core.files TO lodestar_app;
