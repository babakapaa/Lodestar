CREATE SCHEMA IF NOT EXISTS notes;
GRANT USAGE ON SCHEMA notes TO lodestar_app;

CREATE TABLE notes.notes (
	id         uuid        PRIMARY KEY,
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	owner_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	meta       bytea       NOT NULL,
	seq        bigint      NOT NULL DEFAULT 0,
	base_seq   bigint      NOT NULL DEFAULT 0,
	bytes      bigint      NOT NULL DEFAULT 0,
	created_at timestamptz NOT NULL DEFAULT now(),
	changed_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE INDEX notes_owner ON notes.notes (tenant_id, owner_id);
CREATE INDEX notes_deleted ON notes.notes (deleted_at) WHERE deleted_at IS NOT NULL;

CREATE TABLE notes.members (
	tenant_id uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                      REFERENCES core.tenants (id) ON DELETE CASCADE,
	note_id   uuid        NOT NULL REFERENCES notes.notes (id) ON DELETE CASCADE,
	user_id   uuid        NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	can_edit  boolean     NOT NULL,
	added_by  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	added_at  timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (note_id, user_id)
);

CREATE INDEX members_user ON notes.members (tenant_id, user_id);

CREATE TABLE notes.updates (
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	note_id    uuid        NOT NULL REFERENCES notes.notes (id) ON DELETE CASCADE,
	seq        bigint      NOT NULL,
	client_key text        NOT NULL CHECK (length(client_key) BETWEEN 1 AND 100),
	author_id  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	data       bytea       NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	PRIMARY KEY (note_id, seq),
	UNIQUE (note_id, client_key)
);

CREATE TABLE notes.versions (
	id         uuid        PRIMARY KEY DEFAULT gen_random_uuid(),
	tenant_id  uuid        NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                       REFERENCES core.tenants (id) ON DELETE CASCADE,
	note_id    uuid        NOT NULL REFERENCES notes.notes (id) ON DELETE CASCADE,
	seq        bigint      NOT NULL,
	author_id  uuid        REFERENCES core.users (id) ON DELETE SET NULL,
	data       bytea       NOT NULL,
	body       bytea       NOT NULL,
	bytes      bigint      NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	UNIQUE (note_id, seq)
);

DO $$
DECLARE t text;
BEGIN
	FOREACH t IN ARRAY ARRAY['notes', 'members', 'updates', 'versions'] LOOP
		EXECUTE format('ALTER TABLE notes.%I ENABLE ROW LEVEL SECURITY', t);
		EXECUTE format('ALTER TABLE notes.%I FORCE ROW LEVEL SECURITY', t);
		EXECUTE format($p$CREATE POLICY tenant_isolation ON notes.%I
			USING (tenant_id = current_setting('lodestar.tenant')::uuid)
			WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid)$p$, t);
		EXECUTE format('GRANT SELECT, INSERT, UPDATE, DELETE ON notes.%I TO lodestar_app', t);
	END LOOP;
END $$;
