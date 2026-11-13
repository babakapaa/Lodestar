ALTER TABLE sandbox.runs DROP CONSTRAINT runs_kind_check;
ALTER TABLE sandbox.runs ADD CONSTRAINT runs_kind_check CHECK (kind IN ('jupyter', 'vscode', 'job', 'demo'));

ALTER TABLE sandbox.templates ADD COLUMN starter boolean NOT NULL DEFAULT false;
