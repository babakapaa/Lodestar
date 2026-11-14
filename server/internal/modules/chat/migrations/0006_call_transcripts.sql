ALTER TABLE chat.call_recordings
	ADD COLUMN transcript       text CHECK (transcript IN ('pending', 'done', 'failed')),
	ADD COLUMN transcript_tries integer NOT NULL DEFAULT 0,
	ADD COLUMN transcript_at    timestamptz;

CREATE INDEX call_recordings_transcripts ON chat.call_recordings (tenant_id) WHERE transcript = 'pending';
