ALTER TABLE study.assignments ADD COLUMN defense text CHECK (defense IN ('full', 'short', 'sample'));

CREATE TABLE study.defense_slots (
	id            uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	assignment_id uuid        NOT NULL REFERENCES study.assignments (id) ON DELETE CASCADE,
	starts_at     timestamptz NOT NULL,
	minutes       int         NOT NULL CHECK (minutes BETWEEN 5 AND 240),
	place         text        NOT NULL DEFAULT '' CHECK (length(place) <= 500),
	capacity      int         NOT NULL DEFAULT 1 CHECK (capacity BETWEEN 1 AND 50),
	created_by    uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	created_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX defense_slots_assignment ON study.defense_slots (tenant_id, assignment_id, starts_at);

CREATE TABLE study.defenses (
	tenant_id     uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                          REFERENCES core.tenants (id) ON DELETE CASCADE,
	assignment_id uuid        NOT NULL REFERENCES study.assignments (id) ON DELETE CASCADE,
	student_id    uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	slot_id       uuid        REFERENCES study.defense_slots (id) ON DELETE SET NULL,
	selected      boolean     NOT NULL DEFAULT false,
	status        text        NOT NULL DEFAULT 'planned' CHECK (status IN ('planned', 'held', 'missed')),
	questions     jsonb       NOT NULL DEFAULT '[]',
	summary       text        NOT NULL DEFAULT '' CHECK (length(summary) <= 20000),
	examiner_id   uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	held_at       timestamptz,
	updated_at    timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (assignment_id, student_id)
);

CREATE INDEX defenses_slot ON study.defenses (tenant_id, slot_id);
CREATE INDEX defenses_student ON study.defenses (tenant_id, student_id);

ALTER TABLE study.defense_slots ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.defense_slots FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.defense_slots
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

ALTER TABLE study.defenses ENABLE ROW LEVEL SECURITY;
ALTER TABLE study.defenses FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON study.defenses
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, UPDATE, DELETE ON study.defense_slots TO lodestar_app;
GRANT SELECT, INSERT, UPDATE ON study.defenses TO lodestar_app;
