CREATE TABLE team.proposals (
	id          uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	team_id     uuid        NOT NULL REFERENCES team.teams (id) ON DELETE CASCADE,
	proposed_by uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	work        jsonb       NOT NULL,
	note        text        NOT NULL DEFAULT '' CHECK (length(note) <= 1000),
	status      text        NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'sent', 'replaced')),
	silent_at   timestamptz NOT NULL,
	due_at      timestamptz,
	created_at  timestamptz NOT NULL DEFAULT now(),
	sent_at     timestamptz
);

CREATE UNIQUE INDEX proposals_open ON team.proposals (team_id) WHERE status = 'open';
CREATE INDEX proposals_waiting ON team.proposals (silent_at) WHERE status = 'open';

CREATE TABLE team.votes (
	tenant_id   uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                        REFERENCES core.tenants (id) ON DELETE CASCADE,
	proposal_id uuid        NOT NULL REFERENCES team.proposals (id) ON DELETE CASCADE,
	user_id     uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	agreed      boolean     NOT NULL,
	comment     text        NOT NULL DEFAULT '' CHECK (length(comment) <= 1000),
	at          timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (proposal_id, user_id)
);

CREATE TABLE team.cards (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	user_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	interests  text        NOT NULL DEFAULT '' CHECK (length(interests) <= 500),
	skills     text        NOT NULL DEFAULT '' CHECK (length(skills) <= 500),
	hours      text        NOT NULL DEFAULT '' CHECK (length(hours) <= 200),
	role       text        NOT NULL DEFAULT '' CHECK (length(role) <= 100),
	updated_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (tenant_id, user_id)
);

ALTER TABLE team.proposals ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.proposals FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.proposals
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.votes ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.votes FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.votes
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE team.cards ENABLE ROW LEVEL SECURITY;
ALTER TABLE team.cards FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON team.cards
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON team.proposals TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON team.votes TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON team.cards TO lodestar_app;
