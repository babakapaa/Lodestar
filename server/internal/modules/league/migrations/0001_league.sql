CREATE SCHEMA IF NOT EXISTS league;
GRANT USAGE ON SCHEMA league TO lodestar_app;

CREATE TABLE league.leagues (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	origin_kind text        NOT NULL CHECK (length(origin_kind) BETWEEN 1 AND 60),
	origin_id   text        NOT NULL CHECK (length(origin_id) BETWEEN 1 AND 100),
	divisions   jsonb       NOT NULL,
	confirm     boolean     NOT NULL DEFAULT false,
	created_at  timestamptz NOT NULL DEFAULT now(),
	updated_at  timestamptz NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, origin_kind, origin_id)
);

CREATE TABLE league.weights (
	tenant_id uuid NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	               REFERENCES core.tenants (id) ON DELETE CASCADE,
	league_id uuid NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	work_kind text NOT NULL CHECK (length(work_kind) BETWEEN 1 AND 60),
	work_id   text NOT NULL CHECK (length(work_id) BETWEEN 1 AND 100),
	weight    int  NOT NULL CHECK (weight BETWEEN 0 AND 10000),
	PRIMARY KEY (league_id, work_kind, work_id)
);

CREATE TABLE league.seasons (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	league_id  uuid        NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	number     int         NOT NULL CHECK (number >= 1),
	started_at timestamptz NOT NULL DEFAULT now(),
	ended_at   timestamptz,
	UNIQUE (league_id, number)
);

CREATE UNIQUE INDEX seasons_open ON league.seasons (league_id) WHERE ended_at IS NULL;

CREATE TABLE league.players (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	league_id uuid        NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	division  int         NOT NULL DEFAULT 0 CHECK (division >= 0),
	display   text        NOT NULL DEFAULT 'name' CHECK (display IN ('name', 'hidden', 'nickname')),
	nickname  text        NOT NULL DEFAULT '' CHECK (length(nickname) <= 40),
	joined_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (league_id, user_id),
	CHECK (display <> 'nickname' OR nickname <> '')
);

CREATE INDEX players_user ON league.players (tenant_id, user_id);

CREATE TABLE league.entries (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	league_id  uuid        NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	season_id  uuid        NOT NULL REFERENCES league.seasons (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	kind       text        NOT NULL CHECK (kind IN ('work', 'adjustment')),
	work_kind  text        CHECK (length(work_kind) BETWEEN 1 AND 60),
	work_id    text        CHECK (length(work_id) BETWEEN 1 AND 100),
	title      text        NOT NULL DEFAULT '' CHECK (length(title) <= 300),
	score      float8,
	max        float8,
	weight     int,
	points     int         NOT NULL,
	reason     text        NOT NULL DEFAULT '' CHECK (length(reason) <= 1000),
	created_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	CHECK ((kind = 'work') = (work_kind IS NOT NULL AND work_id IS NOT NULL AND score IS NOT NULL AND max IS NOT NULL AND weight IS NOT NULL)),
	CHECK (kind = 'work' OR reason <> '')
);

CREATE UNIQUE INDEX entries_work ON league.entries (league_id, user_id, work_kind, work_id) WHERE kind = 'work';
CREATE INDEX entries_season ON league.entries (season_id, user_id);

CREATE TABLE league.transitions (
	id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	league_id     uuid        NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	season_id     uuid        NOT NULL REFERENCES league.seasons (id) ON DELETE CASCADE,
	user_id       uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	from_division int         NOT NULL CHECK (from_division >= 0),
	to_division   int         NOT NULL CHECK (to_division >= 0 AND to_division <> from_division),
	points        int         NOT NULL,
	status        text        NOT NULL CHECK (status IN ('pending', 'applied', 'rejected', 'reverted')),
	decided_by    uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	decided_at    timestamptz,
	created_at    timestamptz NOT NULL DEFAULT clock_timestamp()
);

CREATE UNIQUE INDEX transitions_pending ON league.transitions (league_id, user_id) WHERE status = 'pending';
CREATE INDEX transitions_season ON league.transitions (season_id, user_id);

CREATE TABLE league.disputes (
	id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	league_id     uuid        NOT NULL REFERENCES league.leagues (id) ON DELETE CASCADE,
	user_id       uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	entry_id      uuid        REFERENCES league.entries (id) ON DELETE CASCADE,
	transition_id uuid        REFERENCES league.transitions (id) ON DELETE CASCADE,
	text          text        NOT NULL CHECK (length(text) BETWEEN 1 AND 2000),
	status        text        NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'accepted', 'rejected')),
	answer        text        NOT NULL DEFAULT '' CHECK (length(answer) <= 2000),
	decided_by    uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	decided_at    timestamptz,
	created_at    timestamptz NOT NULL DEFAULT now(),
	CHECK ((entry_id IS NULL) <> (transition_id IS NULL))
);

CREATE UNIQUE INDEX disputes_open_entry ON league.disputes (entry_id) WHERE status = 'open';
CREATE UNIQUE INDEX disputes_open_transition ON league.disputes (transition_id) WHERE status = 'open';
CREATE INDEX disputes_league ON league.disputes (league_id, status);

ALTER TABLE league.leagues ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.leagues FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.leagues
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.weights ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.weights FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.weights
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.seasons ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.seasons FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.seasons
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.players ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.players FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.players
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.entries ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.entries FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.entries
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.transitions ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.transitions FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.transitions
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE league.disputes ENABLE ROW LEVEL SECURITY;
ALTER TABLE league.disputes FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON league.disputes
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON league.leagues TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON league.weights TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON league.seasons TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON league.players TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON league.entries TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON league.transitions TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON league.disputes TO lodestar_app;
