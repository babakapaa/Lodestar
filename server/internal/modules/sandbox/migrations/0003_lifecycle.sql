ALTER TABLE sandbox.runs DROP CONSTRAINT runs_status_check;
ALTER TABLE sandbox.runs ADD CONSTRAINT runs_status_check
	CHECK (status IN ('queued', 'starting', 'running', 'held', 'stopping', 'finished', 'failed', 'stopped'));

ALTER TABLE sandbox.runs
	ADD COLUMN active_at   timestamptz,
	ADD COLUMN net_bytes   bigint           NOT NULL DEFAULT 0,
	ADD COLUMN idle_warned boolean          NOT NULL DEFAULT false,
	ADD COLUMN gpu_sum     double precision NOT NULL DEFAULT 0,
	ADD COLUMN gpu_samples int              NOT NULL DEFAULT 0,
	ADD COLUMN held_at     timestamptz,
	ADD COLUMN held_secs   double precision NOT NULL DEFAULT 0,
	ADD COLUMN hold_note   text             NOT NULL DEFAULT '' CHECK (length(hold_note) <= 300),
	ADD COLUMN refunded    boolean          NOT NULL DEFAULT false,
	ADD COLUMN log_tail    text             NOT NULL DEFAULT '' CHECK (length(log_tail) <= 8000);

DROP INDEX sandbox.runs_active;
DROP INDEX sandbox.runs_node;
CREATE INDEX runs_active ON sandbox.runs (tenant_id, class_id, status) WHERE status IN ('queued', 'starting', 'running', 'held', 'stopping');
CREATE INDEX runs_node ON sandbox.runs (node_id) WHERE status IN ('starting', 'running', 'held', 'stopping');
CREATE INDEX runs_held ON sandbox.runs (tenant_id, held_at) WHERE status = 'held';
