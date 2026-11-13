package chat

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
)

const Name = "chat"

type Module struct{}

func (Module) Manifest() module.Manifest {
	return module.Manifest{Name: Name, Version: "0.1.0", Optional: true}
}

func (Module) Mount(mux *http.ServeMux, k module.Kernel) {
	Routes{Chats: NewService(k), Messages: NewMessageService(k)}.Mount(mux)
}

func (Module) Workers(k module.Kernel) []module.Worker {
	return []module.Worker{NewRosterWorker(k)}
}

func (Module) Permissions() []string {
	return []string{string(CreateChannels)}
}

func (Module) PrivacySections(k module.Kernel) []privacy.Section {
	return []privacy.Section{PrivacySection{messages: MessageStore{sealer: Sealer{box: k.Secrets}}}}
}
