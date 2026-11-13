package navigator

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/core/module"
	"github.com/babakapaa/lodestar/server/internal/platform/llm"
)

const titleLength = 80

type Service struct {
	k        module.Kernel
	dialogs  DialogStore
	messages MessageStore
}

func NewService(k module.Kernel) Service {
	sealer := Sealer{box: k.Secrets}
	return Service{k: k, dialogs: DialogStore{sealer: sealer}, messages: MessageStore{sealer: sealer}}
}

func (s Service) Dialogs(ctx context.Context, p auth.Principal, scopeKind, scopeID string) ([]Dialog, error) {
	var list []Dialog
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		var err error
		list, err = s.dialogs.List(ctx, tx, p.UserID, scopeKind, scopeID)
		return err
	})
	return list, err
}

func (s Service) Create(ctx context.Context, p auth.Principal, in DialogInput) (DialogDetails, bool, error) {
	in, err := in.normalized()
	if err != nil {
		return DialogDetails{}, false, err
	}
	var d DialogDetails
	var created bool
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		brief, err := s.brief(ctx, tx, p, in.ScopeKind, in.ScopeID)
		if err != nil {
			return err
		}
		dialog, isNew, err := s.dialogs.Insert(ctx, tx, p.UserID, in)
		if err != nil {
			return err
		}
		messages, err := s.messages.List(ctx, tx, dialog.ID)
		d, created = DialogDetails{Dialog: dialog, Brief: brief, Messages: messages}, isNew
		return err
	})
	return d, created, err
}

func (s Service) Get(ctx context.Context, p auth.Principal, id string) (DialogDetails, error) {
	var d DialogDetails
	err := s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		dialog, err := s.dialogs.Owned(ctx, tx, p.UserID, id, false)
		if err != nil {
			return err
		}
		brief, err := s.brief(ctx, tx, p, dialog.ScopeKind, dialog.ScopeID)
		if errors.Is(err, module.ErrBriefNotFound) {
			brief, err = nil, nil
		}
		if err != nil {
			return err
		}
		messages, err := s.messages.List(ctx, tx, dialog.ID)
		d = DialogDetails{Dialog: dialog, Brief: brief, Messages: messages}
		return err
	})
	return d, err
}

func (s Service) Delete(ctx context.Context, p auth.Principal, id string) error {
	return s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		return s.dialogs.Delete(ctx, tx, p.UserID, id)
	})
}

func (s Service) Ask(ctx context.Context, p auth.Principal, id string, in MessageInput, delta func(string) error) (Answer, bool, error) {
	in, err := in.normalized()
	if err != nil {
		return Answer{}, false, err
	}
	if s.k.Model == nil {
		return Answer{}, false, ErrUnavailable
	}
	var prompt Prompt
	var support string
	var done *Answer
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		dialog, err := s.dialogs.Owned(ctx, tx, p.UserID, id, false)
		if err != nil {
			return err
		}
		if answer, ok, err := s.messages.Answered(ctx, tx, id, in.Key); err != nil || ok {
			done = &answer
			return err
		}
		prompt = Prompt{Question: in.Text, Intent: in.Intent, Now: time.Now()}
		if in.More {
			last, err := s.messages.Last(ctx, tx, id)
			if err != nil {
				return err
			}
			prompt.Intent, prompt.Step = last.Intent, last.Step.next()
		}
		if prompt.Brief, err = s.brief(ctx, tx, p, dialog.ScopeKind, dialog.ScopeID); err != nil && !errors.Is(err, module.ErrBriefNotFound) {
			return err
		}
		limit, err := s.memoryLimit(ctx, tx)
		if err != nil {
			return err
		}
		if prompt.Memory, err = s.messages.Memory(ctx, tx, p.UserID, dialog, limit); err != nil {
			return err
		}
		support, err = s.text(ctx, tx, SupportSetting)
		return err
	})
	if err != nil || done != nil {
		if done != nil {
			return *done, false, err
		}
		return Answer{}, false, err
	}
	distress := false
	if prompt.Intent == "" {
		c := Classifier{model: s.k.Model}.Classify(ctx, in.Text)
		prompt.Intent, distress = c.Intent, c.Distress
	}
	if !in.More {
		prompt.Step = prompt.Intent.start()
	}
	reply, err := s.complete(ctx, prompt, delta)
	if err != nil {
		return Answer{}, false, err
	}
	body := messageBody{Text: reply, Intent: prompt.Intent, Step: prompt.Step}
	if distress {
		body.Support = support
	}
	var answer Answer
	var created bool
	err = s.k.Tenants.Run(ctx, p.TenantID, func(tx pgx.Tx) error {
		if _, err := s.dialogs.Owned(ctx, tx, p.UserID, id, true); err != nil {
			return err
		}
		if existing, ok, err := s.messages.Answered(ctx, tx, id, in.Key); err != nil || ok {
			answer = existing
			return err
		}
		question, err := s.messages.Insert(ctx, tx, id, in.Key, RoleStudent, messageBody{Text: in.Text, Intent: prompt.Intent})
		if err != nil {
			return err
		}
		replied, err := s.messages.Insert(ctx, tx, id, "", RoleNavigator, body)
		if err != nil {
			return err
		}
		if _, err := s.dialogs.Touch(ctx, tx, id, title(in.Text)); err != nil {
			return err
		}
		answer, created = Answer{Question: question, Reply: replied}, true
		return nil
	})
	return answer, created, err
}

func (s Service) complete(ctx context.Context, prompt Prompt, delta func(string) error) (string, error) {
	req := llm.Request{Messages: prompt.Messages(), Temperature: 0.3}
	var reply string
	var err error
	if delta != nil {
		reply, err = s.k.Model.Stream(ctx, req, delta)
	} else {
		reply, err = s.k.Model.Complete(ctx, req)
	}
	if err != nil {
		if errors.Is(err, llm.ErrModelFailed) {
			return "", errors.Join(ErrModelFailed, err)
		}
		return "", err
	}
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return "", ErrModelFailed
	}
	return reply, nil
}

func (s Service) brief(ctx context.Context, tx pgx.Tx, p auth.Principal, kind, id string) (*module.Brief, error) {
	if kind == "" {
		return nil, nil
	}
	b, err := s.k.Briefs.Brief(ctx, tx, p, kind, id)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (s Service) memoryLimit(ctx context.Context, tx pgx.Tx) (int, error) {
	raw, err := s.k.Settings.Value(ctx, tx, MemorySetting)
	if err != nil {
		return 0, err
	}
	var n int
	if err := json.Unmarshal(raw, &n); err != nil {
		return 0, err
	}
	return n, nil
}

func (s Service) text(ctx context.Context, tx pgx.Tx, key string) (string, error) {
	raw, err := s.k.Settings.Value(ctx, tx, key)
	if err != nil {
		return "", err
	}
	var v string
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", err
	}
	return strings.TrimSpace(v), nil
}

func title(text string) string {
	text = strings.Join(strings.Fields(text), " ")
	if utf8.RuneCountInString(text) <= titleLength {
		return text
	}
	return string([]rune(text)[:titleLength-1]) + "…"
}
