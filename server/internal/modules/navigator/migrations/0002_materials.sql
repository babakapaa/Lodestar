CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE navigator.materials (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	course_id  uuid        NOT NULL,
	author_id  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	title      text        NOT NULL CHECK (length(title) BETWEEN 1 AND 300),
	reference  text        NOT NULL DEFAULT '' CHECK (length(reference) <= 2000),
	url        text        NOT NULL DEFAULT '' CHECK (length(url) <= 2000),
	created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX materials_course ON navigator.materials (tenant_id, course_id, created_at);

CREATE TABLE navigator.passages (
	tenant_id   uuid   NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                   REFERENCES core.tenants (id) ON DELETE CASCADE,
	material_id uuid   NOT NULL REFERENCES navigator.materials (id) ON DELETE CASCADE,
	seq         int    NOT NULL,
	text        text   NOT NULL,
	model       text,
	embedding   vector,
	PRIMARY KEY (material_id, seq)
);

CREATE INDEX passages_material ON navigator.passages (tenant_id, material_id);

ALTER TABLE navigator.materials ENABLE ROW LEVEL SECURITY;
ALTER TABLE navigator.materials FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON navigator.materials
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE navigator.passages ENABLE ROW LEVEL SECURITY;
ALTER TABLE navigator.passages FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON navigator.passages
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON navigator.materials TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON navigator.passages TO lodestar_app;
