package chat

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/files"
)

type AttachmentStore struct {
	files files.FileStore
}

func (a AttachmentStore) Check(ctx context.Context, tx pgx.Tx, ownerID string, ids []string) error {
	seen := make(map[string]bool, len(ids))
	for _, id := range ids {
		f, err := a.files.ByID(ctx, tx, id, false)
		if err != nil || f.OwnerID != ownerID || f.Status != files.StatusReady || seen[id] {
			return ErrBadFile
		}
		seen[id] = true
	}
	return nil
}

func (AttachmentStore) Attach(ctx context.Context, tx pgx.Tx, messageID string, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO chat.message_files (message_id, file_id, position)
		SELECT $1, f, p FROM unnest($2::uuid[]) WITH ORDINALITY AS x (f, p)`, messageID, ids)
	if err != nil {
		return fmt.Errorf("attach files: %w", err)
	}
	return nil
}

func (AttachmentStore) Detach(ctx context.Context, tx pgx.Tx, messageIDs []string) error {
	_, err := tx.Exec(ctx, `DELETE FROM chat.message_files WHERE message_id::text = ANY($1)`, messageIDs)
	if err != nil {
		return fmt.Errorf("detach files: %w", err)
	}
	return nil
}

func (a AttachmentStore) Fill(ctx context.Context, tx pgx.Tx, list []Message) error {
	if len(list) == 0 {
		return nil
	}
	index := make(map[string]int, len(list))
	ids := make([]string, len(list))
	for i, m := range list {
		index[m.ID], ids[i] = i, m.ID
	}
	rows, err := tx.Query(ctx, `
		SELECT message_id::text, file_id::text FROM chat.message_files
		WHERE message_id::text = ANY($1) ORDER BY message_id, position`, ids)
	if err != nil {
		return fmt.Errorf("message files: %w", err)
	}
	type pair struct{ message, file string }
	pairs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (pair, error) {
		var p pair
		err := row.Scan(&p.message, &p.file)
		return p, err
	})
	if err != nil {
		return fmt.Errorf("message files: %w", err)
	}
	for _, p := range pairs {
		f, err := a.files.ByID(ctx, tx, p.file, false)
		if err != nil {
			return err
		}
		i := index[p.message]
		list[i].Files = append(list[i].Files, f)
	}
	return nil
}
