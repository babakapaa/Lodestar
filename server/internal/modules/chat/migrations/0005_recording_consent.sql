ALTER TABLE chat.calls
	ADD COLUMN consent_requested_by uuid REFERENCES core.users (id) ON DELETE SET NULL,
	ADD COLUMN consent_audio_only   boolean NOT NULL DEFAULT false;

ALTER TABLE chat.call_members ADD COLUMN consent boolean NOT NULL DEFAULT false;

ALTER TABLE chat.call_recordings ADD COLUMN audio_only boolean NOT NULL DEFAULT false;
