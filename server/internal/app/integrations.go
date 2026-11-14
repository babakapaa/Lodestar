package app

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/delivery"
	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/platform/livekit"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

type Integrations struct {
	Storage  files.Storage
	Pusher   delivery.Pusher
	Mailer   delivery.Mailer
	Model    llm.Model
	Embedder llm.Embedder
	Reranker llm.Reranker
	Speech   llm.Transcriber
	Calls    *livekit.Issuer

	LoginClient *http.Client
}
