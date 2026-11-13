ALTER TABLE league.entries DROP CONSTRAINT entries_kind_check;
ALTER TABLE league.entries DROP CONSTRAINT entries_check;
ALTER TABLE league.entries DROP CONSTRAINT entries_check1;
ALTER TABLE league.entries ADD CONSTRAINT entries_kind_check CHECK (kind IN ('work', 'adjustment', 'task'));
ALTER TABLE league.entries ADD CONSTRAINT entries_source_check
	CHECK ((kind = 'adjustment') = (work_kind IS NULL OR work_id IS NULL OR score IS NULL OR max IS NULL OR weight IS NULL));
ALTER TABLE league.entries ADD CONSTRAINT entries_adjustment_check CHECK (kind <> 'adjustment' OR reason <> '');
DROP INDEX league.entries_work;
CREATE UNIQUE INDEX entries_work ON league.entries (league_id, user_id, work_kind, work_id) WHERE kind <> 'adjustment';

CREATE TABLE league.tasks (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	league_id  uuid        NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	title      text        NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
	about      text        NOT NULL DEFAULT '' CHECK (length(about) <= 10000),
	partner    text        NOT NULL DEFAULT '' CHECK (length(partner) <= 200),
	points     int         NOT NULL CHECK (points BETWEEN 1 AND 10000),
	min_size   int         NOT NULL CHECK (min_size BETWEEN 2 AND 12),
	max_size   int         NOT NULL CHECK (max_size BETWEEN min_size AND 12),
	self_formed boolean    NOT NULL DEFAULT true,
	due_at     timestamptz,
	closed_at  timestamptz,
	created_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX tasks_league ON league.tasks (league_id, created_at);

CREATE TABLE league.solutions (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	task_id     uuid        NOT NULL REFERENCES league.tasks (id) ON DELETE CASCADE,
	team_id     text        NOT NULL CHECK (length(team_id) BETWEEN 1 AND 100),
	team_name   text        NOT NULL DEFAULT '' CHECK (length(team_name) <= 200),
	proposal_id text        NOT NULL CHECK (length(proposal_id) BETWEEN 1 AND 100),
	proposed_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	members     uuid[]      NOT NULL,
	work        jsonb       NOT NULL,
	votes       jsonb       NOT NULL DEFAULT '[]',
	embedding   real[],
	similar_to  uuid        REFERENCES league.solutions (id) ON DELETE SET NULL,
	similarity  float8,
	explanation text        NOT NULL DEFAULT '' CHECK (length(explanation) <= 2000),
	points      int,
	passed      uuid[]      NOT NULL DEFAULT '{}',
	comment     text        NOT NULL DEFAULT '' CHECK (length(comment) <= 2000),
	reviewed_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	reviewed_at timestamptz,
	created_at  timestamptz NOT NULL DEFAULT now(),
	UNIQUE (task_id, proposal_id)
);

CREATE INDEX solutions_task ON league.solutions (task_id, team_id, created_at);

CREATE TABLE league.awards (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	title       text        NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
	about       text        NOT NULL DEFAULT '' CHECK (length(about) <= 2000),
	gpu_hours   int         NOT NULL DEFAULT 0 CHECK (gpu_hours BETWEEN 0 AND 100000),
	quantity    int         CHECK (quantity BETWEEN 1 AND 100000),
	archived_at timestamptz,
	created_by  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE league.grants (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	award_id   uuid        NOT NULL REFERENCES league.awards (id) ON DELETE CASCADE,
	league_id  uuid        NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	season_id  uuid        NOT NULL REFERENCES league.seasons (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	reason     text        NOT NULL CHECK (length(reason) BETWEEN 1 AND 1000),
	granted_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	granted_at timestamptz NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX grants_award ON league.grants (award_id);
CREATE INDEX grants_league ON league.grants (league_id, user_id);

ALTER TABLE league.tasks ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.tasks FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.tasks
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.solutions ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.solutions FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.solutions
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.awards ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.awards FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.awards
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.grants ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.grants FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.grants
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON league.tasks TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON league.solutions TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON league.awards TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON league.grants TO lodestar_app;
