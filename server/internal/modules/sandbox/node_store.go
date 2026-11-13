package sandbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const nodeTimeout = 90 * time.Second

const nodeColumns = `n.id::text, n.name, n.class_id::text, c.name, n.gpus,
	(SELECT count(*)::int FROM sandbox.runs r WHERE r.node_id = n.id AND r.status IN ('starting', 'running', 'held', 'stopping')),
	n.seen_at, coalesce(n.seen_at > now() - make_interval(secs => $1), false), n.created_at`

const nodeFrom = ` FROM sandbox.nodes n JOIN sandbox.gpu_classes c ON c.id = n.class_id`

type NodeStore struct{}

func (NodeStore) Insert(ctx context.Context, tx pgx.Tx, in NodeInput, hash []byte) (Node, error) {
	var id string
	err := tx.QueryRow(ctx, `
		INSERT INTO sandbox.nodes (name, class_id, gpus, token_hash) VALUES ($1, $2::uuid, $3, $4) RETURNING id::text`,
		in.Name, in.ClassID, in.GPUs, hash).Scan(&id)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return Node{}, ErrNodeTaken
	}
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return Node{}, ErrClassNotFound
	}
	if err != nil {
		return Node{}, fmt.Errorf("create node: %w", err)
	}
	return NodeStore{}.one(ctx, tx, `n.id::text = $2`, id)
}

func (NodeStore) All(ctx context.Context, tx pgx.Tx) ([]Node, error) {
	rows, err := tx.Query(ctx, `SELECT `+nodeColumns+nodeFrom+` ORDER BY n.name`, nodeTimeout.Seconds())
	if err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	return pgx.CollectRows(rows, scanNode)
}

func (NodeStore) ByToken(ctx context.Context, tx pgx.Tx, hash []byte) (Node, error) {
	n, err := NodeStore{}.one(ctx, tx, `n.token_hash = $2`, hash)
	if errors.Is(err, ErrNodeNotFound) {
		return n, ErrBadNodeToken
	}
	return n, err
}

func (NodeStore) Seen(ctx context.Context, tx pgx.Tx, id string) error {
	if _, err := tx.Exec(ctx, `UPDATE sandbox.nodes SET seen_at = now() WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("mark node seen: %w", err)
	}
	return nil
}

func (NodeStore) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	n, err := NodeStore{}.one(ctx, tx, `n.id::text = $2`, id)
	if err != nil {
		return err
	}
	if n.Busy > 0 {
		return ErrNodeBusy
	}
	if _, err := tx.Exec(ctx, `DELETE FROM sandbox.nodes WHERE id::text = $1`, id); err != nil {
		return fmt.Errorf("delete node: %w", err)
	}
	return nil
}

func (NodeStore) Free(ctx context.Context, tx pgx.Tx, classID string) ([]Node, error) {
	rows, err := tx.Query(ctx, `SELECT `+nodeColumns+nodeFrom+`
		WHERE n.class_id::text = $2 AND n.seen_at > now() - make_interval(secs => $1)
		ORDER BY n.name FOR UPDATE OF n`, nodeTimeout.Seconds(), classID)
	if err != nil {
		return nil, fmt.Errorf("free nodes: %w", err)
	}
	nodes, err := pgx.CollectRows(rows, scanNode)
	if err != nil {
		return nil, err
	}
	free := nodes[:0]
	for _, n := range nodes {
		if n.Busy < n.GPUs {
			free = append(free, n)
		}
	}
	return free, nil
}

func (NodeStore) Lost(ctx context.Context, tx pgx.Tx) ([]Node, error) {
	rows, err := tx.Query(ctx, `SELECT `+nodeColumns+nodeFrom+`
		WHERE coalesce(n.seen_at <= now() - make_interval(secs => $1), true)
		  AND EXISTS (SELECT 1 FROM sandbox.runs r WHERE r.node_id = n.id AND r.status IN ('starting', 'running', 'held', 'stopping'))`,
		nodeTimeout.Seconds())
	if err != nil {
		return nil, fmt.Errorf("lost nodes: %w", err)
	}
	return pgx.CollectRows(rows, scanNode)
}

func (NodeStore) one(ctx context.Context, tx pgx.Tx, where string, arg any) (Node, error) {
	rows, err := tx.Query(ctx, `SELECT `+nodeColumns+nodeFrom+` WHERE `+where, nodeTimeout.Seconds(), arg)
	if err != nil {
		return Node{}, fmt.Errorf("load node: %w", err)
	}
	n, err := pgx.CollectExactlyOneRow(rows, scanNode)
	if errors.Is(err, pgx.ErrNoRows) {
		return Node{}, ErrNodeNotFound
	}
	return n, err
}

func scanNode(row pgx.CollectableRow) (Node, error) {
	var n Node
	err := row.Scan(&n.ID, &n.Name, &n.ClassID, &n.ClassName, &n.GPUs, &n.Busy, &n.SeenAt, &n.Online, &n.CreatedAt)
	return n, err
}

func (NodeStore) Load(ctx context.Context, tx pgx.Tx) (map[string]classLoad, error) {
	rows, err := tx.Query(ctx, `
		SELECT c.id::text,
		       coalesce((SELECT sum(greatest(n.gpus - (SELECT count(*) FROM sandbox.runs r
		                    WHERE r.node_id = n.id AND r.status IN ('starting', 'running', 'held', 'stopping')), 0))
		                 FROM sandbox.nodes n WHERE n.class_id = c.id AND n.seen_at > now() - make_interval(secs => $1)), 0)::int,
		       (SELECT count(*) FROM sandbox.runs r WHERE r.class_id = c.id AND r.status = 'queued')::int
		FROM sandbox.gpu_classes c`, nodeTimeout.Seconds())
	if err != nil {
		return nil, fmt.Errorf("class load: %w", err)
	}
	defer rows.Close()
	out := map[string]classLoad{}
	for rows.Next() {
		var id string
		var l classLoad
		if err := rows.Scan(&id, &l.free, &l.queued); err != nil {
			return nil, fmt.Errorf("class load: %w", err)
		}
		out[id] = l
	}
	return out, rows.Err()
}
