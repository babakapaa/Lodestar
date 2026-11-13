CREATE TABLE core.setup_state (
	tenant_id        uuid        PRIMARY KEY DEFAULT current_setting('lodestar.tenant')::uuid
	                             REFERENCES core.tenants (id) ON DELETE CASCADE,
	protocol_version int         NOT NULL DEFAULT 0,
	signed_by        uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	signed_at        timestamptz,
	launched_by      uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	launched_at      timestamptz,
	started_at       timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE core.setup_sections (
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	section     text        NOT NULL CHECK (length(section) BETWEEN 1 AND 40),
	status      text        NOT NULL DEFAULT 'not_started' CHECK (status IN ('not_started', 'in_progress', 'review', 'approved')),
	assignee_id uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	due_on      date,
	note        text        NOT NULL DEFAULT '' CHECK (length(note) <= 2000),
	updated_by  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	updated_at  timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, section)
);

CREATE INDEX setup_sections_assignee ON core.setup_sections (tenant_id, assignee_id);

CREATE TABLE core.setup_legal (
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	question    text        NOT NULL CHECK (length(question) BETWEEN 1 AND 40),
	answer      text        NOT NULL CHECK (length(answer) BETWEEN 1 AND 4000),
	basis       text        NOT NULL CHECK (length(basis) BETWEEN 1 AND 1000),
	answered_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	answered_at timestamptz NOT NULL DEFAULT now(),
	approved_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	approved_at timestamptz,
	PRIMARY KEY (tenant_id, question)
);

ALTER TABLE core.setup_state ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.setup_state FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.setup_state
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.setup_sections ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.setup_sections FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.setup_sections
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE core.setup_legal ENABLE ROW LEVEL SECURITY;
ALTER TABLE core.setup_legal FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON core.setup_legal
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON core.setup_state TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON core.setup_sections TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON core.setup_legal TO lodestar_app;
