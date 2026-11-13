ALTER TABLE sandbox.ledger DROP CONSTRAINT ledger_kind_check;
ALTER TABLE sandbox.ledger ADD CONSTRAINT ledger_kind_check CHECK (kind IN ('topup', 'transfer', 'run'));
ALTER TABLE sandbox.ledger ADD COLUMN run_id uuid;
CREATE UNIQUE INDEX ledger_run ON sandbox.ledger (run_id) WHERE run_id IS NOT NULL;

CREATE TABLE sandbox.templates (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	name        text        NOT NULL CHECK (length(name) BETWEEN 1 AND 100),
	about       text        NOT NULL DEFAULT '' CHECK (length(about) <= 1000),
	image       text        NOT NULL CHECK (length(image) BETWEEN 1 AND 300),
	archived_at timestamptz,
	created_at  timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, name)
);

CREATE TABLE sandbox.nodes (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	name       text        NOT NULL CHECK (length(name) BETWEEN 1 AND 100),
	class_id   uuid        NOT NULL REFERENCES sandbox.gpu_classes (id),
	gpus       int         NOT NULL CHECK (gpus BETWEEN 1 AND 64),
	token_hash bytea       NOT NULL UNIQUE,
	seen_at    timestamptz,
	created_at timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, name)
);

CREATE TABLE sandbox.runs (
	id         uuid          PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid          NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	term_id    uuid          NOT NULL REFERENCES sandbox.terms (id) ON DELETE CASCADE,
	user_id    uuid          NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	team_id    text          CHECK (length(team_id) BETWEEN 1 AND 100),
	team_name  text          NOT NULL DEFAULT '' CHECK (length(team_name) <= 100),
	kind       text          NOT NULL CHECK (kind IN ('jupyter', 'vscode', 'job')),
	title      text          NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
	command    text          NOT NULL DEFAULT '' CHECK (length(command) <= 2000),
	template_id uuid         NOT NULL REFERENCES sandbox.templates (id),
	class_id   uuid          NOT NULL REFERENCES sandbox.gpu_classes (id),
	rate       numeric(6,2)  NOT NULL,
	max_hours  numeric(6,2)  NOT NULL CHECK (max_hours > 0 AND max_hours <= 168),
	work_kind  text          CHECK (length(work_kind) BETWEEN 1 AND 60),
	work_id    text          CHECK (length(work_id) BETWEEN 1 AND 100),
	work_title text          NOT NULL DEFAULT '' CHECK (length(work_title) <= 300),
	due_at     timestamptz,
	weight     numeric(10,2) NOT NULL DEFAULT 0,
	status     text          NOT NULL DEFAULT 'queued'
	                         CHECK (status IN ('queued', 'starting', 'running', 'stopping', 'finished', 'failed', 'stopped')),
	node_id    uuid          REFERENCES sandbox.nodes (id) ON DELETE SET NULL,
	url        text          NOT NULL DEFAULT '' CHECK (length(url) <= 1000),
	exit_code  int,
	reason     text          NOT NULL DEFAULT '' CHECK (length(reason) <= 1000),
	charged    numeric(10,2),
	queued_at  timestamptz   NOT NULL DEFAULT clock_timestamp(),
	started_at timestamptz,
	ended_at   timestamptz,
	CHECK ((work_kind IS NULL) = (work_id IS NULL))
);

CREATE INDEX runs_active ON sandbox.runs (tenant_id, class_id, status) WHERE status IN ('queued', 'starting', 'running', 'stopping');
CREATE INDEX runs_user ON sandbox.runs (term_id, user_id, queued_at DESC);
CREATE INDEX runs_team ON sandbox.runs (term_id, team_id, queued_at DESC) WHERE team_id IS NOT NULL;
CREATE INDEX runs_node ON sandbox.runs (node_id) WHERE status IN ('starting', 'running', 'stopping');

ALTER TABLE sandbox.templates ENABLE ROW LEVEL SECURITY;
ALTER TABLE sandbox.templates FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sandbox.templates
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE sandbox.nodes ENABLE ROW LEVEL SECURITY;
ALTER TABLE sandbox.nodes FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sandbox.nodes
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE sandbox.runs ENABLE ROW LEVEL SECURITY;
ALTER TABLE sandbox.runs FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sandbox.runs
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON sandbox.templates TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON sandbox.nodes TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON sandbox.runs TO lodestar_app;
