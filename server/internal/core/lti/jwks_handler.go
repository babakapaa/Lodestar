package lti

import (
	"net/http"

	"github.com/go-jose/go-jose/v4"
	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/tenant"
	"github.com/babakapaa/lodestar/server/internal/platform/database"
	"github.com/babakapaa/lodestar/server/internal/platform/jsonhttp"
)

type JWKSHandler struct {
	tenants tenant.Store
	tx      database.TenantTx
	keys    KeyStore
}

func (h JWKSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := h.tenants.BySlug(r.Context(), r.PathValue("tenant"))
	if err != nil {
		writeError(w, err)
		return
	}
	var keys []Key
	err = h.tx.Run(r.Context(), t.ID, func(tx pgx.Tx) error {
		var err error
		keys, err = h.keys.All(r.Context(), tx)
		return err
	})
	if err != nil {
		writeError(w, err)
		return
	}
	set := jose.JSONWebKeySet{Keys: make([]jose.JSONWebKey, 0, len(keys))}
	for _, k := range keys {
		set.Keys = append(set.Keys, jose.JSONWebKey{Key: &k.Private.PublicKey, KeyID: k.ID, Algorithm: string(jose.RS256), Use: "sig"})
	}
	w.Header().Set("Cache-Control", "public, max-age=3600")
	jsonhttp.Write(w, http.StatusOK, set)
}
