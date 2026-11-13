package chat

import (
	"net/http"

	"github.com/babakapaa/lodestar/server/internal/core/files"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/core/privacy"
	"github.com/babakapaa/lodestar/server/internal/core/search"
)

const Name = "chat"

type Module struct{}

func (Module) Manifest() module.Manifest {
	return module.Manifest{Name: Name, Version: "0.1.0", Optional: true}
}

func (Module) Mount(mux *http.ServeMux, k module.Kernel) {
	Routes{Chats: NewService(k), Messages: NewMessageService(k), Blocks: NewBlockService(k), Reports: NewReportService(k)}.Mount(mux)
}

func (Module) Workers(k module.Kernel) []module.Worker {
	return []module.Worker{NewRosterWorker(k), NewPurger(k)}
}

func (Module) Permissions() []string {
	return []string{string(CreateChannels), string(ReviewReports)}
}

func (Module) PrivacySections(k module.Kernel) []privacy.Section {
	sealer := Sealer{box: k.Secrets}
	return []privacy.Section{PrivacySection{messages: MessageStore{sealer: sealer}}, ReportsPrivacySection{reports: ReportStore{sealer: sealer}}}
}

func (Module) SearchSources() []search.Source {
	return []search.Source{ChatsSource{}}
}

func (Module) FileLinks(module.Kernel) []files.Link {
	return []files.Link{FileLink{}}
}
