package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type recipient struct {
	language       string
	tenantLanguage string
	quiet          quietHours
}

type recipientReader struct {
	preferences Preferences
	settings    Settings
}

func (r recipientReader) read(ctx context.Context, tx pgx.Tx, userID string) (recipient, error) {
	var language, tenantLanguage, zone, tenantZone string
	var quietOn bool
	var from, to int
	reads := []struct {
		value any
		load  func() (json.RawMessage, error)
	}{
		{&language, func() (json.RawMessage, error) { return r.preferences.Value(ctx, tx, userID, LanguagePreference) }},
		{&tenantLanguage, func() (json.RawMessage, error) { return r.settings.Value(ctx, tx, LanguageSetting) }},
		{&quietOn, func() (json.RawMessage, error) { return r.preferences.Value(ctx, tx, userID, QuietHoursPreference) }},
		{&from, func() (json.RawMessage, error) { return r.preferences.Value(ctx, tx, userID, QuietFromPreference) }},
		{&to, func() (json.RawMessage, error) { return r.preferences.Value(ctx, tx, userID, QuietToPreference) }},
		{&zone, func() (json.RawMessage, error) { return r.preferences.Value(ctx, tx, userID, TimeZonePreference) }},
		{&tenantZone, func() (json.RawMessage, error) { return r.settings.Value(ctx, tx, TimeZoneSetting) }},
	}
	for _, read := range reads {
		raw, err := read.load()
		if err != nil {
			return recipient{}, err
		}
		if err := json.Unmarshal(raw, read.value); err != nil {
			return recipient{}, fmt.Errorf("read recipient settings: %w", err)
		}
	}
	if language == "system" {
		language = ""
	}
	if tenantLanguage == "system" {
		tenantLanguage = ""
	}
	if zone == "" {
		zone = tenantZone
	}
	location, err := time.LoadLocation(zone)
	if err != nil {
		location = time.UTC
	}
	return recipient{language: language, tenantLanguage: tenantLanguage, quiet: quietHours{on: quietOn, from: from, to: to, location: location}}, nil
}
