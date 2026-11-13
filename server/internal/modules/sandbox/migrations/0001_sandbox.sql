CREATE SCHEMA IF NOT EXISTS sandbox;
GRANT USAGE ON SCHEMA sandbox TO lodestar_app;

CREATE TABLE sandbox.gpu_classes (
	id          uuid         PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id   uuid         NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	name        text         NOT NULL CHECK (length(name) BETWEEN 1 AND 100),
	memory_gb   int          NOT NULL CHECK (memory_gb BETWEEN 1 AND 4096),
	rate        numeric(6,2) NOT NULL CHECK (rate > 0 AND rate <= 1000),
	archived_at timestamptz,
	created_at  timestamptz  NOT NULL DEFAULT now(),
	UNIQUE (tenant_id, name)
);

CREATE TABLE sandbox.terms (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	name       text        NOT NULL CHECK (length(name) BETWEEN 1 AND 100),
	starts_on  date        NOT NULL,
	ends_on    date        NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	CHECK (ends_on > starts_on)
);

CREATE INDEX terms_dates ON sandbox.terms (tenant_id, starts_on);

CREATE TABLE sandbox.topups (
	id         uuid          PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid          NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	term_id    uuid          NOT NULL REFERENCES sandbox.terms (id) ON DELETE CASCADE,
	user_id    uuid          NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	work_kind  text          NOT NULL CHECK (length(work_kind) BETWEEN 1 AND 60),
	work_id    text          NOT NULL CHECK (length(work_id) BETWEEN 1 AND 100),
	title      text          NOT NULL DEFAULT '' CHECK (length(title) <= 300),
	hours      numeric(10,2) NOT NULL CHECK (hours > 0),
	reason     text          NOT NULL DEFAULT '' CHECK (length(reason) <= 1000),
	status     text          NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'declined')),
	granted    numeric(10,2) CHECK (granted > 0),
	decided_by uuid          REFERENCES core.users (id) ON DELETE SET NULL,
	decided_at timestamptz,
	created_at timestamptz   NOT NULL DEFAULT clock_timestamp(),
	CHECK ((status = 'approved') = (granted IS NOT NULL))
);

CREATE UNIQUE INDEX topups_pending ON sandbox.topups (term_id, user_id, work_kind, work_id) WHERE status = 'pending';
CREATE INDEX topups_status ON sandbox.topups (tenant_id, status, created_at);

CREATE TABLE sandbox.ledger (
	id         uuid          PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid          NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                         REFERENCES core.tenants (id) ON DELETE CASCADE,
	term_id    uuid          NOT NULL REFERENCES sandbox.terms (id) ON DELETE CASCADE,
	owner_kind text          NOT NULL CHECK (owner_kind IN ('user', 'team')),
	owner_id   text          NOT NULL CHECK (length(owner_id) BETWEEN 1 AND 100),
	kind       text          NOT NULL CHECK (kind IN ('topup', 'transfer')),
	hours      numeric(10,2) NOT NULL CHECK (hours <> 0),
	title      text          NOT NULL DEFAULT '' CHECK (length(title) <= 300),
	actor_id   uuid          REFERENCES core.users (id) ON DELETE SET NULL,
	created_at timestamptz   NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX ledger_owner ON sandbox.ledger (term_id, owner_kind, owner_id);

ALTER TABLE sandbox.gpu_classes ENABLE ROW LEVEL SECURITY;
ALTER TABLE sandbox.gpu_classes FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sandbox.gpu_classes
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE sandbox.terms ENABLE ROW LEVEL SECURITY;
ALTER TABLE sandbox.terms FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sandbox.terms
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE sandbox.topups ENABLE ROW LEVEL SECURITY;
ALTER TABLE sandbox.topups FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sandbox.topups
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE sandbox.ledger ENABLE ROW LEVEL SECURITY;
ALTER TABLE sandbox.ledger FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON sandbox.ledger
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE ON sandbox.gpu_classes TO lodestar_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON sandbox.terms TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON sandbox.topups TO lodestar_app;
GRANT SELECT, INSERT ON sandbox.ledger TO lodestar_app;
