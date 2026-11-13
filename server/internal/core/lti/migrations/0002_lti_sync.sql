ALTER TABLE lti.contexts
	ADD COLUMN assignment_source text CHECK (assignment_source IN ('lodestar', 'lms')),
	ADD COLUMN synced_at         timestamptz,
	ADD COLUMN sync_error        text;

GRANT DELETE ON lti.members TO lodestar_app;

CREATE TABLE lti.line_items (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	context_id uuid        NOT NULL REFERENCES lti.contexts (id) ON DELETE CASCADE,
	url        text        NOT NULL,
	work_kind  text        NOT NULL,
	work_id    text        NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	UNIQUE (context_id, url),
	UNIQUE (context_id, work_kind, work_id)
);

CREATE TABLE lti.scores (
	line_item_id uuid        NOT NULL REFERENCES lti.line_items (id) ON DELETE CASCADE,
	user_id      uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	tenant_id    uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	score        float8      NOT NULL,
	max          float8      NOT NULL,
	sent_at      timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (line_item_id, user_id)
);

ALTER TABLE lti.line_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.line_items FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.line_items
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE lti.scores ENABLE ROW LEVEL SECURITY;
ALTER TABLE lti.scores FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON lti.scores
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON lti.line_items TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON lti.scores TO lodestar_app;
