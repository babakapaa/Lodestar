package livekit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const recorderTTL = 10 * time.Minute

func (i *Issuer) Records() bool {
	return i.cfg.Storage.Endpoint != "" && i.cfg.Storage.AccessKey != ""
}

func (i *Issuer) StartRecording(ctx context.Context, room, key string, audioOnly bool, now time.Time) (string, error) {
	if !i.Records() {
		return "", ErrNoRecording
	}
	st := i.cfg.Storage
	var out EgressInfo
	err := i.call(ctx, "StartRoomCompositeEgress", map[string]any{
		"room_name":  room,
		"layout":     "grid",
		"audio_only": audioOnly,
		"file_outputs": []map[string]any{{
			"file_type": fileType(audioOnly),
			"filepath":  key,
			"s3": map[string]any{
				"access_key":       st.AccessKey,
				"secret":           st.SecretKey,
				"region":           st.Region,
				"endpoint":         st.Endpoint,
				"bucket":           st.Bucket,
				"force_path_style": true,
			},
		}},
	}, &out, now)
	if err != nil {
		return "", err
	}
	if out.EgressID == "" {
		return "", ErrRecordingRefused
	}
	return out.EgressID, nil
}

func (i *Issuer) StopRecording(ctx context.Context, egressID string, now time.Time) error {
	return i.call(ctx, "StopEgress", map[string]any{"egress_id": egressID}, nil, now)
}

func (i *Issuer) call(ctx context.Context, method string, in, out any, now time.Time) error {
	token, err := i.sign("lodestar", map[string]any{"video": videoClaims{RoomRecord: true}}, recorderTTL, now)
	if err != nil {
		return err
	}
	body, err := json.Marshal(in)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.cfg.api()+"/twirp/livekit.Egress/"+method, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := i.http.Do(req)
	if err != nil {
		return fmt.Errorf("livekit %s: %w", method, err)
	}
	defer func() { _ = resp.Body.Close() }()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("livekit %s: %w", method, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s %d", ErrRecordingRefused, method, resp.StatusCode)
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(raw, out)
}
