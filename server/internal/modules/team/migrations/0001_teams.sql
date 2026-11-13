CREATE SCHEMA IF NOT EXISTS team;
GRANT USAGE ON SCHEMA team TO lodestar_app;

CREATE TABLE team.teams (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	scope_kind  text        NOT NULL CHECK (length(scope_kind) BETWEEN 1 AND 60),
	scope_id    text        NOT NULL CHECK (length(scope_id) BETWEEN 1 AND 100),
	name        text        NOT NULL CHECK (length(name) BETWEEN 1 AND 100),
	invite_code text        NOT NULL UNIQUE CHECK (length(invite_code) BETWEEN 16 AND 64),
	created_by  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	created_at  timestamptz NOT NULL DEFAULT now(),
	archived_at timestamptz
);

CREATE INDEX teams_scope ON team.teams (tenant_id, scope_kind, scope_id);

CREATE TABLE team.members (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	team_id    uuid        NOT NULL REFERENCES team.teams (id) ON DELETE CASCADE,
	scope_kind text        NOT NULL,
	scope_id   text        NOT NULL,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	label      text        NOT NULL DEFAULT '' CHECK (length(label) <= 40),
	joined_at  timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (team_id, user_id),
	UNIQUE (tenant_id, scope_kind, scope_id, user_id)
);

CREATE INDEX members_user ON team.members (tenant_id, user_id);

CREATE TABLE team.history (
	id        uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	team_id   uuid        NOT NULL REFERENCES team.teams (id) ON DELETE CASCADE,
	user_id   uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	event     text        NOT NULL CHECK (event IN ('created', 'joined', 'added', 'left', 'removed', 'separated', 'moved_in', 'moved_out')),
	actor_id  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	at        timestamptz NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX history_team ON team.history (team_id, at);

CREATE TABLE team.invitations (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	team_id    uuid        NOT NULL REFERENCES team.teams (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	invited_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (team_id, user_id)
);

CREATE INDEX invitations_user ON team.invitations (tenant_id, user_id);

CREATE TABLE team.seekers (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	scope_kind text        NOT NULL,
	scope_id   text        NOT NULL,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	note       text        NOT NULL DEFAULT '' CHECK (length(note) <= 500),
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, scope_kind, scope_id, user_id)
);

CREATE TABLE team.solos (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	scope_kind text        NOT NULL,
	scope_id   text        NOT NULL,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	decided_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, scope_kind, scope_id, user_id)
);

CREATE TABLE team.leave_requests (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	team_id    uuid        NOT NULL REFERENCES team.teams (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	reason     text        NOT NULL DEFAULT '' CHECK (length(reason) <= 1000),
	status     text        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
	created_at timestamptz NOT NULL DEFAULT now(),
	decided_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	decided_at timestamptz
);

CREATE UNIQUE INDEX leave_requests_pending ON team.leave_requests (team_id, user_id) WHERE status = 'pending';

ALTER TABLE team.teams ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.teams FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.teams
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.members ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.members FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.members
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.history ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.history FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.history
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.invitations ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.invitations FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.invitations
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.seekers ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.seekers FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.seekers
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.solos ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.solos FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.solos
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.leave_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.leave_requests FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.leave_requests
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON team.teams TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON team.members TO lodestar_app;
GRANT SELECT, INSERT ON team.history TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON team.invitations TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON team.seekers TO lodestar_app;
GRANT SELECT, INSERT, DELETE ON team.solos TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON team.leave_requests TO lodestar_app;
