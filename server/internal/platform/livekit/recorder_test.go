package livekit_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"

	"github.com/babakapaa/lodestar/server/internal/platform/livekit"
	"github.com/babakapaa/lodestar/server/internal/platform/objectstore"
)

func TestStartRecordingSendsRoomAndStorage(t *testing.T) {
	now := time.Now()
	var got map[string]any
	var path string
	var record bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		parsed, err := jwt.ParseSigned(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "), []jose.SignatureAlgorithm{jose.HS256})
		if err != nil {
			t.Errorf("parse: %v", err)
		}
		var extra struct {
			Video struct {
				Record bool `json:"roomRecord"`
			} `json:"video"`
		}
		if err := parsed.Claims([]byte(secret), &extra); err != nil {
			t.Errorf("claims: %v", err)
		}
		record = extra.Video.Record
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &got)
		_, _ = w.Write([]byte(`{"egress_id":"EG_1","egressId":"EG_1"}`))
	}))
	defer srv.Close()
	issuer := livekit.NewIssuer(livekit.Config{
		URL: "ws" + strings.TrimPrefix(srv.URL, "http"), APIKey: "key", APISecret: secret,
		Storage: objectstore.Config{Endpoint: "http://minio:9000", Bucket: "lodestar", AccessKey: "a", SecretKey: "s"},
	})
	id, err := issuer.StartRecording(context.Background(), "t_c", "t/f", now)
	if err != nil || id != "EG_1" {
		t.Fatalf("start = %q, %v", id, err)
	}
	if path != "/twirp/livekit.Egress/StartRoomCompositeEgress" || !record || got["room_name"] != "t_c" {
		t.Fatalf("request = %s %v %v", path, record, got)
	}
	out := got["file_outputs"].([]any)[0].(map[string]any)
	if out["filepath"] != "t/f" || out["s3"].(map[string]any)["bucket"] != "lodestar" {
		t.Fatalf("output = %v", out)
	}
}

func TestRecordingNeedsStorage(t *testing.T) {
	issuer := livekit.NewIssuer(livekit.Config{URL: "wss://calls", APIKey: "key", APISecret: secret})
	if _, err := issuer.StartRecording(context.Background(), "t_c", "t/f", time.Now()); !errors.Is(err, livekit.ErrNoRecording) {
		t.Fatalf("err = %v", err)
	}
}

func TestEgressEventReadsSizeAsString(t *testing.T) {
	var e livekit.Event
	raw := `{"event":"egress_ended","egressInfo":{"egressId":"EG_1","roomName":"t_c","status":"EGRESS_COMPLETE","fileResults":[{"filename":"t/f","size":"2048"}]}}`
	if err := json.Unmarshal([]byte(raw), &e); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if e.EgressInfo == nil || e.EgressInfo.FileResults[0].Size != 2048 {
		t.Fatalf("event = %+v", e.EgressInfo)
	}
}
