CREATE TABLE navigator.topic_hits (
	tenant_id   uuid NOT NULL DEFAULT current_setting('lodestar.tenant')::uuid
	                 REFERENCES core.tenants (id) ON DELETE CASCADE,
	course_id   text NOT NULL CHECK (length(course_id) BETWEEN 1 AND 100),
	material_id uuid REFERENCES navigator.materials (id) ON DELETE CASCADE,
	user_id     uuid NOT NULL REFERENCES core.users (id) ON DELETE CASCADE,
	intent      text NOT NULL DEFAULT '' CHECK (length(intent) <= 30),
	asked_on    date NOT NULL
);

CREATE INDEX topic_hits_course ON navigator.topic_hits (tenant_id, course_id, asked_on);

ALTER TABLE navigator.topic_hits ENABLE ROW LEVEL SECURITY;
ALTER TABLE navigator.topic_hits FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON navigator.topic_hits
	USING (tenant_id = current_setting('lodestar.tenant')::uuid)
	WITH CHECK (tenant_id = current_setting('lodestar.tenant')::uuid);

GRANT SELECT, INSERT, DELETE ON navigator.topic_hits TO lodestar_app;
